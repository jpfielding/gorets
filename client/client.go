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

// RetsSession is an interface defining the expected
type RetsSession interface {
	ChangePassword(ctx context.Context, url string) error
	Get(ctx context.Context, url string) error
	GetMetadata(ctx context.Context, r MetadataRequest) (*Metadata, error)
	GetObject(ctx context.Context, r GetObjectRequest) (<-chan GetObjectResult, error)
	GetPayloadList(ctx context.Context, p PayloadListRequest) (*PayloadList, error)
	Login(ctx context.Context, url string) (*CapabilityURLs, error)
	Logout(ctx context.Context, logoutURL string) (*LogoutResponse, error)
	PostObject(ctx context.Context, url string) error
	Search(ctx context.Context, r SearchRequest) (*SearchResult, error)
	Update(ctx context.Context, url string) error
}

// Session holds the state of the server interaction
type Session struct {
	HTTPMethodDefault string

	Username, Password string

	Execute func(ctx context.Context, req *http.Request) (*http.Response, error)

	Reset func() error
}

// OnRequest ...
type OnRequest func(req *http.Request)

// NewSession configures the default rets session
func NewSession(user, pw, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (*Session, error) {
	var session Session

	session.HTTPMethodDefault = "GET"

	session.Username = user
	session.Password = pw

	if transport == nil {
		transport = http.DefaultTransport
	}

	client := http.Client{
		Transport: transport,
	}

	session.Reset = func() error {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}
		client.Jar = jar
		return nil
	}
	err := session.Reset()
	if err != nil {
		return nil, err
	}

	var onRequest []OnRequest
	// apply default headers
	onRequest = append(onRequest, func(req *http.Request) {
		req.Header.Set(UserAgent, userAgent)
		req.Header.Set(RETSVersion, retsVersion)
		req.Header.Set(Accept, "*/*")
	})
	// apply ua auth headers per request
	if userAgentPw != "" {
		uaAuth := UserAgentAuthentication{
			UserAgent:         userAgent,
			UserAgentPassword: userAgentPw,
		}
		onRequest = append(onRequest, uaAuth.OnRequest)
	}

	session.Execute = func(ctx context.Context, req *http.Request) (*http.Response, error) {
		for _, or := range onRequest {
			or(req)
		}
		return ctxhttp.Do(ctx, &client, req)
	}

	return &session, nil
}
