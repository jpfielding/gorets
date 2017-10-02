package proxy

import (
	"fmt"
	"net/http"
	"strings"
)

// Config TODO does any auth require the version info for auth
type Config struct {
	URL                          string
	User, Password               string
	UserAgent, UserAgentPassword string
}

// Login manages de/multiplexing requests to RETS servers
func Login(prefix string, proxy map[string][]Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, prefix)
		parts := strings.Split(sub, "/")
		src := parts[0]
		usr := parts[1]
		fmt.Fprintf(res, "Hi there, I love proxying %s for %s login!", src, usr)
	}
}
