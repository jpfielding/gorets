//
// Provides the basic mechanism for User Agent authentication for rets
//
// RETS 1.8 - 3.10 Computing the RETS-UA-Authorization Value

package client

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

// UserAgentAuthentication ...
type UserAgentAuthentication struct {
	transport http.RoundTripper

	RETSVersion string

	UserAgent,
	UserAgentPassword string

	// TODO consider a hook for providing a RETS-Request-ID
}

// RoundTrip meets the http.RoundTripper interface to apply UAuth
func (t *UserAgentAuthentication) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// did someone else set this?
	requestId := req.Header.Get(RETS_REQUEST_ID)
	sessionId := ""
	if h, err := req.Cookie(RETS_SESSION_ID); err == nil {
		sessionId = h.Value
	}
	uaAuthHeader := calculateUaAuthHeader(
		t.UserAgent,
		t.UserAgentPassword,
		requestId,
		sessionId,
		t.RETSVersion,
	)
	// this will replace an existing value
	req.Header.Set(RETS_UA_AUTH_HEADER, uaAuthHeader)
	return t.transport.RoundTrip(req)
}

func calculateUaAuthHeader(userAgent, userAgentPw, requestId, sessionId, retsVersion string) string {
	hasher := md5.New()

	io.WriteString(hasher, userAgent+":"+userAgentPw)
	secretHash := hex.EncodeToString(hasher.Sum(nil))

	pieces := strings.Join([]string{secretHash, requestId, sessionId, retsVersion}, ":")

	hasher.Reset()
	io.WriteString(hasher, pieces)
	return "Digest " + hex.EncodeToString(hasher.Sum(nil))
}
