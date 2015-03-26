/**
provides a wrapper and built in auth at the transport
layer.
*/
package client

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
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
	ChangePassword(url string) error
	Get(url string) error
	GetMetadata(r MetadataRequest) (*Metadata, error)
	GetObject(quit <-chan struct{}, r GetObjectRequest) (<-chan GetObjectResult, error)
	GetPayloadList(p PayloadListRequest) (*PayloadList, error)
	Login(url string) (*CapabilityUrls, error)
	Logout(logoutUrl string) (*LogoutResponse, error)
	PostObject(url string) error
	Search(r SearchRequest, quit <-chan struct{}) (*SearchResult, error)
	Update(url string) error
}

/* holds the state of the server interaction */
type Session struct {
	Username, Password string

	UserAgent, UserAgentPassword string

	Version string
	Accept  string

	Cookies []*http.Cookie

	Client http.Client
}

func NewSession(user, pw, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (*Session, error) {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.7"
	session.UserAgent = userAgent
	session.UserAgentPassword = userAgentPw
	session.Version = retsVersion
	session.Accept = "*/*"
	session.Cookies = make([]*http.Cookie, 0)

	if transport == nil {
		transport = http.DefaultTransport
	}

	retsTransport := RetsTransport{
		transport: transport,
		session:   session,
		digest:    nil,
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	session.Client = http.Client{
		Transport: &retsTransport,
		Jar:       jar,
	}
	return &session, nil
}

/* wrapper to intercept each http call */
type RetsTransport struct {
	transport http.RoundTripper
	session   Session
	digest    *Digest
}

func (t *RetsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add(USER_AGENT, t.session.UserAgent)
	req.Header.Add(RETS_VERSION, t.session.Version)
	req.Header.Add(ACCEPT, t.session.Accept)
	for _, cookie := range t.session.Cookies {
		req.AddCookie(cookie)
	}

	if t.session.UserAgentPassword != "" {
		requestId := req.Header.Get(RETS_REQUEST_ID)
		sessionId := ""
		if h, err := req.Cookie(RETS_SESSION_ID); err == nil {
			sessionId = h.Value
		}
		uaAuthHeader := CalculateUaAuthHeader(
			t.session.UserAgent,
			t.session.UserAgentPassword,
			requestId,
			sessionId,
			t.session.Version,
		)
		req.Header.Add(RETS_UA_AUTH_HEADER, uaAuthHeader)
	}
	res, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return res, err
	}
	// TODO check to see if im going to do anything different, if not, just return
	if err = res.Body.Close(); err != nil {
		return res, err
	}
	if t.digest != nil {
		req.Header.Add(WWW_AUTH_RESP, t.digest.CreateDigestResponse(t.session.Username, t.session.Password, req.Method, req.URL.Path))
	}
	t.session.Cookies = make([]*http.Cookie, 0)
	for _, cookie := range res.Cookies() {
		t.session.Cookies = append(t.session.Cookies, cookie)
		req.AddCookie(cookie)
	}
	challenge := res.Header.Get(WWW_AUTH)

	if strings.HasPrefix(strings.ToLower(challenge), "basic") {
		req.SetBasicAuth(t.session.Username, t.session.Password)
		return t.transport.RoundTrip(req)
	} else if strings.HasPrefix(strings.ToLower(challenge), "digest") {
		t.digest, err = NewDigest(challenge)
		if err != nil {
			return nil, err
		}
		req.Header.Set(WWW_AUTH_RESP, t.digest.CreateDigestResponse(t.session.Username, t.session.Password, req.Method, req.URL.Path))
		return t.transport.RoundTrip(req)
	}
	return nil, errors.New("unknown authentication challenge: " + challenge)
}
