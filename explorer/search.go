package explorer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// SearchArgs ...
type SearchArgs struct {
	ID     string       `json:"id"`
	Params SearchParams `json:"params"`
}

// SearchParams ...
type SearchParams struct {
	Resource  string `json:"resource"`
	Class     string `json:"class"`
	Format    string `json:"format"`
	QueryType string `json:"query-type"`
	CountType int    `json:"count-type"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Query     string `json:"query"`
	Select    string `json:"select"`
}

// SearchPage ...
type SearchPage struct {
	Columns rets.Row   `json:"columns"`
	Rows    []rets.Row `json:"rows"`
	MaxRows bool       `json:"maxrows"`
}

// SearchService ...
type SearchService struct{}

// Run ....
func (ms SearchService) Run(r *http.Request, args *SearchArgs, reply *SearchPage) error {
	fmt.Printf("search run params: %v\n", args)
	if args.Params.QueryType == "" {
		args.Params.QueryType = "DQML2"
	}
	if args.Params.Format == "" {
		args.Params.Format = "COMPACT_DECODED"
	}
	s := sessions.Open(args.ID)
	fmt.Printf("%v\n", s.Connection)
	ctx := context.Background()
	return s.Exec(ctx, func(r rets.Requester, u rets.CapabilityURLs, err error) error {
		req := rets.SearchRequest{
			URL: u.Search,
			SearchParams: rets.SearchParams{
				Select:     args.Params.Select,
				Query:      args.Params.Query,
				SearchType: args.Params.Resource,
				Class:      args.Params.Class,
				Format:     args.Params.Format,
				QueryType:  args.Params.QueryType,
				Count:      args.Params.CountType,
				Limit:      args.Params.Limit,
				Offset:     args.Params.Offset,
			},
		}
		fmt.Printf("Querying : %v\n", req)
		if err != nil {
			return err
		}
		result, err := rets.SearchCompact(r, ctx, req)
		defer result.Close()
		if err != nil {
			return nil
		}
		// opening the strea
		reply.Columns = result.Columns
		// too late to err in http here, need another solution
		reply.MaxRows, err = result.ForEach(func(row rets.Row, err error) error {
			reply.Rows = append(reply.Rows, row)
			return err
		})
		return err
	})
}
