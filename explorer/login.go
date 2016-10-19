package explorer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// LoginParams info to use the demo
type LoginParams struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserAgent   string `json:"user-agent"`
	UserAgentPw string `json:"user-agent-pw"`
	Version     string `json:"rets-version"`
}

// Login ...
// input: LoginParams
// output: rets.CapabilityURLS
func Login(ctx context.Context, c *Connection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p LoginParams
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Printf("params: %v\n", p)
		// start with the default Dialer from http.DefaultTransport
		transport := wirelog.NewHTTPTransport()
		// logging
		if c.WireLogFile != "" {
			err = wirelog.LogToFile(transport, c.WireLogFile, true, true)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			fmt.Println("wire logging enabled:", c.WireLogFile)
		}
		requester, err := rets.DefaultSession(
			p.Username,
			p.Password,
			p.UserAgent,
			p.UserAgentPw,
			p.Version,
			transport,
		)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		urls, err := rets.Login(requester, ctx, rets.LoginRequest{URL: p.URL})
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(*urls)
		c.URLs = *urls
		c.Requester = requester
	}
}
