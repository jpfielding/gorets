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

// GetObject ...
func GetObject(ops map[string]string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, ops["GetObject"])
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

		// TODO also check body in case of POST params in body
		values := req.URL.Query()
		params := rets.GetObjectParams{
			Resource: values.Get("Resource"),
			Type:     values.Get("Type"),
			ID:       values.Get("ID"),
			UID:      values.Get("UID"),
		}

		if l := values.Get("Location"); l != "" {
			params.Location, _ = strconv.Atoi(l)
		}

		ctx := context.Background()
		response, err := rets.GetObjects(ctx, r, rets.GetObjectRequest{
			URL:                   urls.GetObject,
			HTTPMethod:            req.Method,
			HTTPFormEncodedValues: (req.Method == "POST"),
			GetObjectParams:       params,
		})
		defer response.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(res, "get objects err %s", err)
			return
		}
		// success, send the urls (modified to point to this server)
		res.Header()["Content-Type"] = response.Header["Content-Type"]
		res.WriteHeader(http.StatusOK)
		io.Copy(res, response.Body)
	}
}
