package explorer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// SearchArgs ...
type SearchArgs struct {
	ID        string `json:"id"`
	Resource  string `json:"resource"`
	Class     string `json:"class"`
	Format    string `json:"format"`
	Select    string `json:"select"`
	CountType int    `json:"count-type"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Query     string `json:"query"`
	QueryType string `json:"query-type"`
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
		args.QueryType = "DQML2"
	}
	if args.Format == "" {
		args.Format = "COMPACT_DECODED"
	}
	s := sessions.Open(args.ID)
	if s == nil {
		return fmt.Errorf("no source found for %s", args.ID)
	}
	fmt.Printf("%v\n", s.Connection)
	ctx := context.Background()
	return s.Exec(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
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
		reply.Count = result.Count
		return err
	})
}
