package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// LoginRequest info to use the demo
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
func Login(s Session) func(http.ResponseWriter, *http.Request) {
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
		fmt.Println(p.Username)
		// start with the default Dialer from http.DefaultTransport
		transport := wirelog.NewHTTPTransport()
		// logging
		if s.WireLogFile != "" {
			err = wirelog.LogToFile(transport, s.WireLogFile, true, true)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			xlog.Println("wire logging enabled:", c.WireLogFile)
		}
		req, err := gorets.DefaultSession(
			p.Username,
			p.Password,
			p.UserAgent,
			p.UserAgentPw,
			p.RetsVersion,
			transport,
		)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		// TODO deal with contexts in the web appropriately
		ctx := context.Background()
		urls, err := rets.Login(s.Requester, ctx, p.URL)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.URLs = urls
		s.session = req
		json.NewEncoder(w).Encode(*urls)
	}
}
