package client

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

// UserAgentAuthentication ...
// RETS 1.8 - 3.10 Computing the RETS-UA-Authorization Value
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
	requestID := req.Header.Get(RETSRequestID)
	sessionID := ""
	if h, err := req.Cookie(RETSSessionID); err == nil {
		sessionID = h.Value
	}
	uaAuthHeader := calculateUaAuthHeader(
		t.UserAgent,
		t.UserAgentPassword,
		requestID,
		sessionID,
		t.RETSVersion,
	)
	// this will replace an existing value
	req.Header.Set(RETSUAAuth, uaAuthHeader)
	return t.transport.RoundTrip(req)
}

func calculateUaAuthHeader(userAgent, userAgentPw, requestID, sessionID, retsVersion string) string {
	hasher := md5.New()

	io.WriteString(hasher, userAgent+":"+userAgentPw)
	secretHash := hex.EncodeToString(hasher.Sum(nil))

	pieces := strings.Join([]string{secretHash, requestID, sessionID, retsVersion}, ":")

	hasher.Reset()
	io.WriteString(hasher, pieces)
	return "Digest " + hex.EncodeToString(hasher.Sum(nil))
}
