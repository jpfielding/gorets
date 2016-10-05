package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
)

func main() {
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
			panic(err)
		}
	}
	fmt.Printf("Connection Settings: %v\n", config)
	if *searchFile != "" {
		err := searchOpts.LoadFrom(*searchFile)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Search Options: %v\n", searchOpts)

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	capability, err := rets.Login(session, ctx, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	defer rets.Logout(session, ctx, rets.LogoutRequest{URL: capability.Logout})

	fmt.Println("Search: ", capability.Search)
	req := rets.SearchRequest{
		URL: capability.Search,
		SearchParams: rets.SearchParams{
			Select:     searchOpts.Select,
			Query:      searchOpts.Query,
			SearchType: searchOpts.Resource,
			Class:      searchOpts.Class,
			Format:     searchOpts.Format,
			QueryType:  searchOpts.QueryType,
			Count:      searchOpts.CountType,
			Limit:      searchOpts.Limit,
			// Offset:     -1,
		},
	}

	w := csv.NewWriter(os.Stdout)
	if *output != "" {
		os.MkdirAll(*output, 0777)
		f, _ := os.Create(*output + "/results.csv")
		defer f.Close()
		w = csv.NewWriter(f)
	}
	defer w.Flush()

	for hasMoreRows := true; hasMoreRows; {
		fmt.Printf("Querying next page: %v\n", req)
		result, err := rets.SearchCompact(session, ctx, req)
		defer result.Close()
		if err != nil {
			panic(err)
		}
		w.Write(result.Columns)
		count := 0
		hasMoreRows, err = result.ForEach(func(row rets.Row, err error) error {
			w.Write(row)
			count++
			return nil
		})
		if hasMoreRows {
			req.Offset = req.Offset + count
		}
	}
}

// SearchOptions ...
type SearchOptions struct {
	Resource  string `json:"resource"`
	Class     string `json:"class"`
	Format    string `json:"format"`
	QueryType string `json:"query-type"`
	CountType int    `json:"count-type"`
	Limit     int    `json:"limit"`
	Query     string `json:"query"`
	Select    string `json:"select"`
}

// SetFlags ...
func (o *SearchOptions) SetFlags() {
	flag.StringVar(&o.Resource, "resource", "Property", "Resource for the search")
	flag.StringVar(&o.Class, "class", "Residential", "Subtype of resource")
	flag.StringVar(&o.Format, "format", "COMPACT-DECODED", "Format for the RETS response")
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
