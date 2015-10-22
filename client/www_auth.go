package client

import (
	"log"
	"net/http"
	"strings"
)

// WWWAuthTransport manages rfc2617 authentication
type WWWAuthTransport struct {
	transport          http.RoundTripper
	Username, Password string
	Digest             *Digest
	HasBasic           bool
}

func (t *WWWAuthTransport) digestResponse(req *http.Request) string {
	return t.Digest.CreateDigestResponse(t.Username, t.Password, req.Method, req.URL.Path)
}

// RoundTrip ...
func (t *WWWAuthTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// attempt to preempt the challenge to save a req/res round trip
	switch {
	case t.Digest != nil:
		req.Header.Set(WWWAuthResp, t.digestResponse(req))
	case t.HasBasic:
		req.SetBasicAuth(t.Username, t.Password)
	}
	res, err := t.transport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	// check for auth issues
	if res.StatusCode != http.StatusUnauthorized {
		return res, err
	}
	// TODO what should we do if we get more than one challenge type?
	for _, c := range res.Header[WWWAuth] {
		switch {
		case strings.HasPrefix(strings.ToLower(c), "digest"):
			t.Digest, err = NewDigest(c)
			if err != nil {
				log.Println("failed to process digest", c, err)
				t.Digest = nil
				continue
			}
			req.Header.Set(WWWAuthResp, t.digestResponse(req))
			return t.transport.RoundTrip(req)
		case strings.HasPrefix(strings.ToLower(c), "basic"):
			t.HasBasic = true
			req.SetBasicAuth(t.Username, t.Password)
			return t.transport.RoundTrip(req)
		}
	}
	return res, err
}
