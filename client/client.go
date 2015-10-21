package client

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/context"
)

// const DefaultTimeout int = 300000

/* standard http header names */
const (
	UserAgent   string = "User-Agent"
	Accept      string = "Accept"
	ContentType string = "Content-Type"
	WWWAuth     string = "Www-Authenticate"
	WWWAuthResp string = "Authorization"
)

/* rets http header names */
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
	Username, Password string

	UserAgent string

	Version string
	Accept  string

	Client http.Client
}

// NewSession configures the default rets session
func NewSession(user, pw, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (*Session, error) {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.7"
	session.UserAgent = userAgent
	session.Version = retsVersion
	session.Accept = "*/*"

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	session.Client = http.Client{
		Transport: transport,
		Jar:       jar,
	}

	if session.Client.Transport == nil {
		session.Client.Transport = http.DefaultTransport
	}

	if userAgentPw != "" {
		session.Client.Transport = &UserAgentAuthentication{
			RETSVersion:       retsVersion,
			UserAgent:         userAgent,
			UserAgentPassword: userAgentPw,
			transport:         session.Client.Transport,
		}
	}

	session.Client.Transport = &WWWAuthTransport{
		transport: session.Client.Transport,
		Username:  user,
		Password:  pw,
	}

	session.Client.Transport = &RetsTransport{
		transport: session.Client.Transport,
		session:   session,
	}

	return &session, nil
}

// RetsTransport to intercept each http call
type RetsTransport struct {
	transport http.RoundTripper
	session   Session
}

// RoundTrip implements http.RoundTripper
func (t *RetsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add(UserAgent, t.session.UserAgent)
	req.Header.Add(RETSVersion, t.session.Version)
	req.Header.Add(Accept, t.session.Accept)

	return t.transport.RoundTrip(req)
}
