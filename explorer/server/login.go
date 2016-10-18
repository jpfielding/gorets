package server

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
func Login(ctx context.Context, u *User) func(http.ResponseWriter, *http.Request) {
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
		if u.WireLogFile != "" {
			err = wirelog.LogToFile(transport, u.WireLogFile, true, true)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			fmt.Println("wire logging enabled:", u.WireLogFile)
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
		// TODO deal with contexts in the web appropriately
		urls, err := rets.Login(u.Requester, ctx, rets.LoginRequest{URL: p.URL})
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(*urls)
		u.URLs = *urls
		u.Requester = requester
	}
}
