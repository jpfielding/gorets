package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jpfielding/gorets/rets"
)

// Search ...
func Search(ops map[string]string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, ops["Search"])
		parts := strings.Split(sub, "/")
		src := parts[0]
		usr := parts[1]
		if _, ok := srcs[src]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "source %s not found", src)
			return
		}
		if _, ok := srcs[src][usr]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "user %s not found", usr)
			return
		}
		session := srcs[src][usr]
		r, urls, err := session.Get()
		if err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(res, "source %s, user %s login failed", src, usr)
			return
		}

		// TOOD also check body in case of POST params in body
		values := req.URL.Query()
		params := rets.SearchParams{
			SearchType: values.Get("SearchType"),
			Class:      values.Get("Class"),
			Query:      values.Get("Query"),
		}
		// TODO apply optionally
		params.Format = values.Get("Format")
		params.Select = values.Get("Select")
		params.Payload = values.Get("Payload")
		params.QueryType = values.Get("QueryType")
		params.RestrictedIndicator = values.Get("RestrictedIndicator")
		if c := values.Get("Count"); c != "" {
			params.Count, _ = strconv.Atoi(c)
		}
		if l := values.Get("Limit"); l != "" {
			if strings.ToUpper(l) == "NONE" {
				params.Limit = -1
			} else {
				params.Limit, _ = strconv.Atoi(l)
			}
		}
		if o := values.Get("Offset"); o != "" {
			params.Offset, _ = strconv.Atoi(o)
		}
		if s := values.Get("StandardNames"); s != "" {
			params.StandardNames, _ = strconv.Atoi(s)
		}
		// TODO standardnames

		// if we're posting
		if req.Method == "POST" {
			params.HTTPFormEncodedValues = true
		}

		ctx := context.Background()
		reader, err := rets.SearchStream(ctx, r, rets.SearchRequest{
			URL:          urls.Search,
			HTTPMethod:   req.Method,
			SearchParams: params,
		})
		defer reader.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(res, "search err %s", err)
			return
		}
		// success, send the urls (modified to point to this server)
		// TODO set content-type here
		res.WriteHeader(http.StatusOK)
		io.Copy(res, reader)
	}
}
