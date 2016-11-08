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
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
)

func main() {
	status := 0
	defer func() { os.Exit(status) }()

	configFile := flag.String("config-file", "", "Config file for RETS connection")
	searchFile := flag.String("search-options", "", "Config file for search options")
	output := flag.String("output", "", "Directory for file output")

	config := common.Config{}
	config.SetFlags()

	searchOpts := SearchOptions{}
	searchOpts.SetFlags()

	flag.Parse()

	if *configFile != "" {
		err := config.LoadFrom(*configFile)
		if err != nil {
			log.Println(err)
			status = 1
			return
		}
	}
	log.Printf("Connection Settings: %v\n", config)
	if *searchFile != "" {
		err := searchOpts.LoadFrom(*searchFile)
		if err != nil {
			log.Println(err)
			status = 2
			return
		}
	}
	log.Printf("Search Options: %v\n", searchOpts)

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		log.Println(err)
		status = 3
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	urls, err := rets.Login(session, ctx, rets.LoginRequest{URL: config.URL})
	if err != nil {
		log.Println(err)
		status = 4
		return
	}
	defer rets.Logout(session, ctx, rets.LogoutRequest{URL: urls.Logout})

	if urls.GetPayloadList != "" {
		log.Println("Payloads: ", urls.GetPayloadList)
		payloads, err := rets.GetPayloadList(session, ctx, rets.PayloadListRequest{
			URL: urls.GetPayloadList,
			ID:  fmt.Sprintf("%s:%s", searchOpts.Resource, searchOpts.Class),
		})
		if err != nil {
			log.Println(err)
			status = 5
			return
		}
		err = payloads.ForEach(func(payload rets.CompactData, err error) error {
			fmt.Printf("%v\n", payload)
			return err
		})
		if err != nil {
			log.Println(err)
			status = 6
			return
		}
	}

	log.Println("Search: ", urls.Search)
	req := rets.SearchRequest{
		URL: urls.Search,
		SearchParams: rets.SearchParams{
			Select:     searchOpts.Select,
			Query:      searchOpts.Query,
			SearchType: searchOpts.Resource,
			Class:      searchOpts.Class,
			Format:     searchOpts.Format,
			QueryType:  searchOpts.QueryType,
			Count:      searchOpts.CountType,
			Limit:      searchOpts.Limit,
		},
	}
	if strings.HasPrefix(req.SearchParams.Format, "COMPACT") {
		processCompact(session, ctx, req, output)
	} else {
		processXML(session, ctx, req, searchOpts.ElementName, output)
	}
}

func processXML(sess rets.Requester, ctx context.Context, req rets.SearchRequest, elem string, output *string) {
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
		log.Printf("Querying next page: %v\n", req)
		result, err := rets.StandardXMLSearch(sess, ctx, req)
		if err != nil {
			log.Println(err)
			return
		}
		switch result.Response.Code {
		case rets.StatusOK:
			// we got some daters
		case rets.StatusNoRecords:
			return
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			log.Println(errors.New(result.Response.Text))
			return
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
			log.Println(err)
			return
		}
		result.Close()
		if err != nil {
			log.Println(err)
			return
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

func processCompact(sess rets.Requester, ctx context.Context, req rets.SearchRequest, output *string) {
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
		log.Printf("Querying next page: %v\n", req)
		result, err := rets.SearchCompact(sess, ctx, req)
		if err != nil {
			log.Println(err)
			return
		}
		switch result.Response.Code {
		case rets.StatusOK:
			// we got some daters
		case rets.StatusNoRecords:
			return
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			log.Println(errors.New(result.Response.Text))
			return
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
			log.Println(err)
			return
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

// SearchOptions ...
type SearchOptions struct {
	Resource string `json:"resource"`
	Class    string `json:"class"`

	Format string `json:"format"`
	Select string `json:"select"`

	Query     string `json:"query"`
	QueryType string `json:"query-type"`

	Payload string `json:"payload"`

	// necessary for xml formats and extracting the sub dom
	ElementName string `json:"element-name"`

	// StandardNames *int   string `json:"standard-names"`
	CountType int `json:"count-type"`
	Limit     int `json:"limit"`
}

// SetFlags ...
func (o *SearchOptions) SetFlags() {
	flag.StringVar(&o.Resource, "resource", "Property", "Resource for the search")
	flag.StringVar(&o.Class, "class", "Residential", "Subtype of resource")
	flag.StringVar(&o.Format, "format", "COMPACT-DECODED", "Format for the RETS response")
	flag.StringVar(&o.Payload, "payload", "", "Requested payload format")
	flag.StringVar(&o.QueryType, "query-type", "DMQL2", "Query type (defaults to DMQL2)")
	flag.StringVar(&o.Select, "select", "", "Fields to be returned")
	flag.StringVar(&o.Query, "dmql", "(ModificationTimestamp=2000-01-01T00:00:00+)", "DMQL for the results")
	flag.IntVar(&o.CountType, "count-type", rets.CountIncluded, "How to deal with the search count")
	flag.IntVar(&o.Limit, "limit", 0, "Limit rows returned per page")
}

// LoadFrom ...
func (o *SearchOptions) LoadFrom(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, o)
	if err != nil {
		return err
	}
	return nil
}
