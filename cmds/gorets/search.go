package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
)

func NewSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use: "search",
		Short: "Search a RETS server",
		Run : search,
	}
}

func search(cmd *cobra.Command, args []string) {
	config := cmd.Flags().GetString("config", "", "Path to the config info for RETS connection")
	output := cmd.Flags().GetString("output", "", "Directory for file output")
	timeout := cmd.Flags().GetInt("timeout", 60, "Seconds to timeout the connection"

	// the params for talking to the server
	connect := Connect{}
	LoadFrom(path.Join(config, "connect.json"), &connect)

	// the params for the query to run
	params := rets.SearchParams{}
	LoadFrom(config, "search.json"), params)

	// config overrides
	resource := cmd.Flags().GetString("resource", params.SearchType, "Resource type for the RETS search, overrides search file")
	class := cmd.Flags().GetString("class", params.Class, "Resource type for the RETS search, overrides search file")
	query := cmd.Flags().GetString("query", params.Query, "Query for the RETS search, overrides search file")

	if dmql != "" {
		params.Query = query
	}
	// launch the session
	session, err := connect.Initialize()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Seconds)
	defer cancel()

	urls, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: urls.Logout})

	fmt.Println("Search: ", urls.Search)
	req := rets.SearchRequest{
		URL: urls.Search,
		SearchParams: params,
	}
	if strings.HasPrefix(req.SearchParams.Format, "COMPACT") {
		processCompact(ctx, session, req, output)
	} else {
		processXML(ctx, session, req, searchOpts.ElementName, output)
	}
}

func processXML(ctx context.Context, sess rets.Requester, req rets.SearchRequest, elem string, output *string) {
	w := os.Stdout
	if *output != "" {
		os.MkdirAll(*output, 0777)
		f, _ := os.Create(*output + "/results.xml")
		w = f
		defer f.Close()
	}

	if elem == "" {
		elem = ".*Listing$"
	}
	rgx, _ := regexp.Compile(elem)
	match := func(t xml.StartElement) bool {
		return rgx.MatchString(t.Name.Local)
	}
	// loop over all the pages we need
	for {
		fmt.Printf("Querying next page: %v\n", req)
		result, err := rets.StandardXMLSearch(ctx, sess, req)
		if err != nil {
			panic(err)
		}
		switch result.Response.Code {
		case rets.StatusOK:
			// we got some daters
		case rets.StatusNoRecords:
			return
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			panic(errors.New(result.Response.Text))
		}
		count := 0
		_, more, err := result.ForEach(match, func(row io.ReadCloser, err error) error {
			if err != nil {
				return err
			}
			count++
			io.Copy(w, row)
			return err
		})
		if err != nil {
			panic(err)
		}
		result.Close()
		if err != nil {
			panic(err)
		}
		if !more {
			return
		}
		if req.Offset == 0 {
			req.Offset = 1
		}
		req.Offset = req.Offset + count
	}
}

func processCompact(ctx context.Context, sess rets.Requester, req rets.SearchRequest, output *string) {
	w := csv.NewWriter(os.Stdout)
	if *output != "" {
		os.MkdirAll(*output, 0777)
		f, _ := os.Create(*output + "/results.csv")
		defer f.Close()
		w = csv.NewWriter(f)
	}
	defer w.Flush()

	// loop over all the pages we need
	for {
		fmt.Printf("Querying next page: %v\n", req)
		result, err := rets.SearchCompact(ctx, sess, req)
		if err != nil {
			panic(err)
		}
		switch result.Response.Code {
		case rets.StatusOK:
			// we got some daters
		case rets.StatusNoRecords:
			return
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			panic(errors.New(result.Response.Text))
		}
		count := 0
		if count == 0 {
			w.Write(result.Columns)
		}
		hasMoreRows, err := result.ForEach(func(row rets.Row, err error) error {
			if err != nil {
				return err
			}
			w.Write(row)
			count++
			return err
		})
		result.Close()
		if err != nil {
			panic(err)
		}
		if !hasMoreRows {
			return
		}
		if req.Offset == 0 {
			req.Offset = 1
		}
		req.Offset = req.Offset + count
	}
}

