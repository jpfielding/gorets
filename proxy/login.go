package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jpfielding/gorets/rets"
)

// Login manages de/multiplexing requests to RETS servers
func Login(prefix string, srcs Sources) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, prefix)
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
		_, urls, err := session.Get()
		if err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(res, "source %s, user %s login failed", src, usr)
			return
		}
		// success, send the urls (modified to point to this server)
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, asXML(*urls))
	}
}

func asXML(urls rets.CapabilityURLs) string {
	return "Hi there, I love proxying for login!"
}
