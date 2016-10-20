package explorer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// SearchParams ...
type SearchParams struct {
	ID        string `json:"id"`
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
func (ms SearchService) Run(r *http.Request, args *SearchParams, reply *SearchPage) error {
	fmt.Printf("search run params: %v\n", args)
	if args.QueryType == "" {
		args.QueryType = "DQML2"
	}
	if args.Format == "" {
		args.Format = "COMPACT_DECODED"
	}
	c := (&ConnectionService{}).Load()[args.ID]
	req := rets.SearchRequest{
		URL: c.URLs.Search,
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
	ctx := context.Background()
	rqr, err := c.Login(ctx)
	if err != nil {
		return err
	}
	result, err := rets.SearchCompact(rqr, ctx, req)
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
}
