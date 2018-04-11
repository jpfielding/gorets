package rets

import (
	"net/http"
	"strings"

	"context"
)

// WWWAuthTransport manages rfc2617 authentication
type WWWAuthTransport struct {
	Requester          Requester
	Username, Password string
	digester           *Digest
	hasBasic           bool
}

func (auth *WWWAuthTransport) header(req *http.Request) string {
	return auth.digester.CreateDigestResponse(
		auth.Username,
		auth.Password,
		req.Method,
		req.URL.Path,
	)
}

// Request ...
func (auth *WWWAuthTransport) Request(ctx context.Context, req *http.Request) (*http.Response, error) {
	// attempt to preempt the challenge to save a req/res round trip
	if auth.hasBasic {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	// apply basic second in case they request basic
	if auth.digester != nil {
		req.Header.Set(WWWAuthResp, auth.header(req))
	}
	res, err := auth.Requester(ctx, req)
	if err != nil {
		return res, err
	}
	// check for auth issues
	if res.StatusCode != http.StatusUnauthorized {
		return res, err
	}
	authed := false
	// sometimes we get more than one challenge type?
	for _, c := range res.Header[WWWAuth] {
		if strings.HasPrefix(strings.ToLower(c), "basic") {
			auth.hasBasic = true
			req.SetBasicAuth(auth.Username, auth.Password)
			authed = true
		}
		if strings.HasPrefix(strings.ToLower(c), "digest") {
			auth.digester, err = NewDigest(c)
			if err != nil {
				auth.digester = nil
				return res, err
			}
			req.Header.Set(WWWAuthResp, auth.header(req))
		}
	}
	if authed {
		return auth.Requester(ctx, req)
	}
	return res, err
}
