package client

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

// RequestIDer allows functions to be provided to generate request ids
type RequestIDer func(req *http.Request) string

// UserAgentAuthentication ...
// RETS 1.8 - 3.10 Computing the RETS-UA-Authorization Value
type UserAgentAuthentication struct {
	UserAgent,
	UserAgentPassword string

	GetRequestID RequestIDer
}

// OnRequest allows ua-auth to be hooked into requests prior to sending
func (ua *UserAgentAuthentication) OnRequest(req *http.Request) {
	// this should already be set
	retsVersion := req.Header.Get(RETSVersion)
	// we generate this and set it in the headers
	requestID := ""
	if ua.GetRequestID != nil {
		requestID = ua.GetRequestID(req)
	}
	sessionID := ""
	if h, err := req.Cookie(RETSSessionID); err == nil {
		sessionID = h.Value
	}
	uaAuthHeader := ua.generateHeader(requestID, sessionID, retsVersion)
	// this will replace an existing value
	req.Header.Set(RETSUAAuth, uaAuthHeader)
}

func (ua *UserAgentAuthentication) generateHeader(requestID, sessionID, version string) string {
	hasher := md5.New()

	io.WriteString(hasher, ua.UserAgent+":"+ua.UserAgentPassword)
	secretHash := hex.EncodeToString(hasher.Sum(nil))

	pieces := strings.Join([]string{secretHash, requestID, sessionID, version}, ":")

	hasher.Reset()
	io.WriteString(hasher, pieces)
	return "Digest " + hex.EncodeToString(hasher.Sum(nil))
}
