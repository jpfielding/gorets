package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jpfielding/gorets/rets"
)

// Metadata ...
func Metadata(ops map[string]string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, ops["GetMetadata"])
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

		ctx := context.Background()
		reader, err := rets.MetadataStream(ctx, r, rets.MetadataRequest{
			URL:    urls.GetMetadata,
			Format: values.Get("Format"),
			MType:  values.Get("Type"),
			ID:     values.Get("ID"),
		})
		defer reader.Close()
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "metadata err %s", err)
			return
		}
		// success, send the urls (modified to point to this server)
		res.WriteHeader(http.StatusOK)
		io.Copy(res, reader)
	}
}
