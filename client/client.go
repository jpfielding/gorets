package client

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// const DefaultTimeout int = 300000

// standard http header names
const (
	UserAgent   string = "User-Agent"
	Accept      string = "Accept"
	ContentType string = "Content-Type"
	WWWAuth     string = "Www-Authenticate"
	WWWAuthResp string = "Authorization"
)

// rets http header names
const (
	RETSVersion   string = "RETS-Version"
	RETSSessionID string = "RETS-Session-ID"
	RETSRequestID string = "RETS-Request-ID"
	RETSUAAuth    string = "RETS-UA-Authorization"
)

// TODO create a Session interface with a Requester and a reset to clear state and pass that in

// Requester ...
type Requester func(ctx context.Context, req *http.Request) (*http.Response, error)

// DefaultSession configures the default rets session
func DefaultSession(user, pwd, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (Requester, error) {
	if transport == nil {
		transport = http.DefaultTransport
	}

	client := http.Client{
		Transport: transport,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	// send the request
	request := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return ctxhttp.Do(ctx, &client, req)
	}
	// www auth
	wwwAuth := (&WWWAuthTransport{
		Requester: request,
		Username:  user,
		Password:  pwd,
	}).Request
	// apply ua auth headers per request, if there is a pwd
	uaAuth := (&UserAgentAuthentication{
		Requester:         wwwAuth,
		UserAgent:         userAgent,
		UserAgentPassword: userAgentPw,
	}).Request
	// apply default headers first (outermost wrapping)
	headers := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		req.Header.Set(UserAgent, userAgent)
		req.Header.Set(RETSVersion, retsVersion)
		req.Header.Set(Accept, "*/*")
		return uaAuth(ctx, req)
	}
	return headers, nil
}
