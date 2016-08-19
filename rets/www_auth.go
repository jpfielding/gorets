package rets

import (
	"log"
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
	switch {
	case auth.digester != nil:
		req.Header.Set(WWWAuthResp, auth.header(req))
	case auth.hasBasic:
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	res, err := auth.Requester(ctx, req)
	if err != nil {
		return res, err
	}
	// check for auth issues
	if res.StatusCode != http.StatusUnauthorized {
		return res, err
	}
	// sometimes we get more than one challenge type?
	for _, c := range res.Header[WWWAuth] {
		switch {
		case strings.HasPrefix(strings.ToLower(c), "digest"):
			auth.digester, err = NewDigest(c)
			if err != nil {
				log.Println("failed to process digest", c, err)
				auth.digester = nil
				continue
			}
			req.Header.Set(WWWAuthResp, auth.header(req))
			return auth.Requester(ctx, req)
		case strings.HasPrefix(strings.ToLower(c), "basic"):
			auth.hasBasic = true
			req.SetBasicAuth(auth.Username, auth.Password)
			return auth.Requester(ctx, req)
		}
	}
	return res, err
}
