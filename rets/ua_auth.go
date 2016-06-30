package rets

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

// UserAgentAuthentication ...
// RETS 1.8 - 3.10 Computing the RETS-UA-Authorization Value
type UserAgentAuthentication struct {
	Requester Requester

	UserAgent,
	UserAgentPassword string
	// usually a constant, but a func hook just in case it changes mid-session
	GetRETSVersion RequestIDer
	// go create one if you want, sadist
	CreateRequestID RequestIDer
	// extracts from a cookie, but that cookie isnt guaranteed to be in the req that we receive
	GetSessionID RequestIDer
}

// Request allows ua-auth to be hooked into requests prior to sending
func (ua *UserAgentAuthentication) Request(ctx context.Context, req *http.Request) (*http.Response, error) {
	// nothing to do gtfo
	if ua.UserAgentPassword == "" {
		return ua.Requester(ctx, req)
	}
	//
	retsVersion := ""
	if ua.GetRETSVersion != nil {
		retsVersion = ua.GetRETSVersion(req)
	}
	// we generate this and set it in the headers
	requestID := ""
	if ua.CreateRequestID != nil {
		requestID = ua.CreateRequestID(req)
		req.Header.Set(RETSRequestID, requestID)
	}
	// need to extract this from context that we may not have at this level
	sessionID := ""
	if ua.GetSessionID != nil {
		sessionID = ua.GetSessionID(req)
	}
	uaAuthHeader := ua.header(requestID, sessionID, retsVersion)
	// this will replace an existing value
	req.Header.Set(RETSUAAuth, uaAuthHeader)
	return ua.Requester(ctx, req)
}

func (ua *UserAgentAuthentication) header(requestID, sessionID, version string) string {
	hasher := md5.New()

	io.WriteString(hasher, ua.UserAgent+":"+ua.UserAgentPassword)
	secretHash := hex.EncodeToString(hasher.Sum(nil))

	pieces := strings.Join([]string{secretHash, requestID, sessionID, version}, ":")
	hasher.Reset()
	io.WriteString(hasher, pieces)
	return "Digest " + hex.EncodeToString(hasher.Sum(nil))
}

// RequestIDer allows functions to be provided to generate request ids
type RequestIDer func(req *http.Request) string

// CreateSessionIDer provides a default implement for extracting the session from a cookie jar
func CreateSessionIDer(jar http.CookieJar) RequestIDer {
	return func(req *http.Request) string {
		for _, c := range jar.Cookies(req.URL) {
			if c.Name == RETSSessionID {
				return c.Value
			}
		}
		return ""
	}
}

// CreateRETSVersioner provides a hook to get a RETS Version from a request
func CreateRETSVersioner(version string) RequestIDer {
	return func(req *http.Request) string {
		return version
	}
}
