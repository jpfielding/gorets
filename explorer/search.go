package explorer

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// SearchArgs ...
type SearchArgs struct {
	Connection Config `json:"connection"`
	Resource   string `json:"resource"`
	Class      string `json:"class"`
	Format     string `json:"format"`
	Select     string `json:"select"`
	CountType  int    `json:"count-type"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Query      string `json:"query"`
	QueryType  string `json:"query-type"`
}

// SearchPage ...
type SearchPage struct {
	Columns rets.Row   `json:"columns"`
	Rows    []rets.Row `json:"rows"`
	MaxRows bool       `json:"maxrows"`
	Count   int        `json:"count"`
}

// SearchService ...
type SearchService struct{}

// Run ....
func (ms SearchService) Run(r *http.Request, args *SearchArgs, reply *SearchPage) error {
	fmt.Printf("search run params: %v\n", args)
	if args.QueryType == "" {
		args.QueryType = "DMQL2"
	}
	if args.Format == "" {
		args.Format = "COMPACT_DECODED"
	}
	cfg := args.Connection
	ctx := context.Background()
	sess, err := cfg.Connect(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", cfg)
	return sess.Process(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		req := rets.SearchRequest{
			URL: u.Search,
			SearchParams: rets.SearchParams{
				Select:     args.Select,
				Query:      args.Query,
				SearchType: args.Resource,
				Class:      args.Class,
				Format:     args.Format,
				QueryType:  args.QueryType,
				Count:      args.CountType,
				Limit:      args.Limit,
				Offset:     args.Offset,
			},
		}
		fmt.Printf("Querying : %v\n", req)
		result, err := rets.SearchCompact(ctx, r, req)
		defer result.Close()
		if err != nil {
			return nil
		}
		// non success rets codes should return an error
		switch result.Response.Code {
		case rets.StatusOK, rets.StatusNoRecords:
		default: // shit hit the fan
			fmt.Printf("Querying : %v\n", result.Response.Text)
			return errors.New(result.Response.Text)
		}
		// opening the strea
		reply.Columns = result.Columns
		// too late to err in http here, need another solution
		reply.MaxRows, err = result.ForEach(func(row rets.Row, err error) error {
			reply.Rows = append(reply.Rows, row)
			return err
		})
		reply.Count = result.Count
		return err
	})
}