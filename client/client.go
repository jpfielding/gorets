/**
provides a wrapper and built in auth at the transport
layer.
*/
package client

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/context"
)

const DEFAULT_TIMEOUT int = 300000

/* standard http header names */
const (
	USER_AGENT    string = "User-Agent"
	ACCEPT        string = "Accept"
	CONTENT_TYPE  string = "Content-Type"
	WWW_AUTH      string = "Www-Authenticate"
	WWW_AUTH_RESP string = "Authorization"
)

/* rets http header names */
const (
	RETS_VERSION        string = "RETS-Version"
	RETS_SESSION_ID     string = "RETS-Session-ID"
	RETS_REQUEST_ID     string = "RETS-Request-ID"
	RETS_UA_AUTH_HEADER string = "RETS-UA-Authorization"
)

type RetsSession interface {
	ChangePassword(ctx context.Context, url string) error
	Get(ctx context.Context, url string) error
	GetMetadata(ctx context.Context, r MetadataRequest) (*Metadata, error)
	GetObject(ctx context.Context, r GetObjectRequest) (<-chan GetObjectResult, error)
	GetPayloadList(ctx context.Context, p PayloadListRequest) (*PayloadList, error)
	Login(ctx context.Context, url string) (*CapabilityUrls, error)
	Logout(ctx context.Context, logoutUrl string) (*LogoutResponse, error)
	PostObject(ctx context.Context, url string) error
	Search(ctx context.Context, r SearchRequest) (*SearchResult, error)
	Update(ctx context.Context, url string) error
}

/* holds the state of the server interaction */
type Session struct {
	Username, Password string

	UserAgent string

	Version string
	Accept  string

	Client http.Client
}

func NewSession(user, pw, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (*Session, error) {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.7"
	session.UserAgent = userAgent
	session.Version = retsVersion
	session.Accept = "*/*"

	if transport == nil {
		transport = http.DefaultTransport
	}

	if userAgentPw != "" {
		transport = &UserAgentAuthentication{
			RETSVersion:       retsVersion,
			UserAgent:         userAgent,
			UserAgentPassword: userAgentPw,
			transport:         transport,
		}
	}

	transport = &WWWAuthTransport{
		transport: transport,
		Username:  user,
		Password:  pw,
	}

	transport = &RetsTransport{
		transport: transport,
		session:   session,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	session.Client = http.Client{
		Transport: transport,
		Jar:       jar,
	}
	return &session, nil
}

/* wrapper to intercept each http call */
type RetsTransport struct {
	transport http.RoundTripper
	session   Session
}

func (t *RetsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add(USER_AGENT, t.session.UserAgent)
	req.Header.Add(RETS_VERSION, t.session.Version)
	req.Header.Add(ACCEPT, t.session.Accept)

	return t.transport.RoundTrip(req)
}
