package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

const (
	searchFileFlag = "params"
)

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search a RETS server",
		Run:   search,
	}
	cmd.Flags().String(searchFileFlag, "", "Search parameters file location")
	return cmd
}

func search(cmd *cobra.Command, args []string) {
	connect, output, timeout := getPersistentFlagValues(search, cmd)

	searchFile, err := cmd.Flags().GetString(searchFileFlag)
	handleError(search, err)

	// the params for the query to run (with some defaults)
	params := rets.SearchParams{
		SearchType: "Property",
		Class:      "Residential",
		Format:     "COMPACT-DECODED",
		Offset:     1,
		QueryType:  "DMQL2",
	}
	err = LoadFrom(searchFile, &params)
	handleError(search, err)

	// should we throw an err here too?
	session, err := connect.Initialize()
	handleError(search, err)

	// launch the session
	session, err = connect.Initialize()
	handleError(search, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	urls, err := rets.Login(ctx, session, rets.LoginRequest{URL: connect.URL})
	handleError(search, err)
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: urls.Logout})

	fmt.Println("Search: ", urls.Search)
	req := rets.SearchRequest{
		URL:          urls.Search,
		SearchParams: params,
	}
	if !strings.HasPrefix(req.SearchParams.Format, "COMPACT") {
		handleError(search, fmt.Errorf("unsupported format %s", params.Format))
	}
	processCompact(ctx, session, req, output)
}

func processCompact(ctx context.Context, sess rets.Requester, req rets.SearchRequest, output string) {
	w := csv.NewWriter(os.Stdout)
	if output != "" {
		os.MkdirAll(output, 0777)
		f, _ := os.Create(output + "/results.csv")
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
