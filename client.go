/**
	provides a wrapper and built in auth at the transport
	layer.
 */
package gorets

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"io"
)

const DEFAULT_TIMEOUT int = 300000

/* standard http header names */
const (
	USER_AGENT string = "User-Agent"
	ACCEPT string = "Accept"
	CONTENT_TYPE string = "Content-Type"
	WWW_AUTH string = "Www-Authenticate"
	WWW_AUTH_RESP string = "Authorization"
)
/* standard http gzip header names */
const (
	ACCEPT_ENCODING string = "Accept-Encoding"
	CONTENT_ENCODING string = "Content-Encoding"
	DEFLATE_ENCODINGS string = "gzip,deflate"
)
/* rets http header names */
const (
	RETS_VERSION string = "RETS-Version"
	RETS_SESSION_ID string = "RETS-Session-ID"
	RETS_REQUEST_ID string = "RETS-Request-ID"
	RETS_UA_AUTH_HEADER string = "RETS-UA-Authorization"
)

/* holds the state of the server interaction */
type Session struct {
	Username,Password string

	UserAgent, UserAgentPassword string
	HttpMethod string

	Version string
	Accept string

	Client http.Client
}


/* the common wrapper details for each response */
type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

func NewSession(user, pw, userAgent, userAgentPw string, logger io.WriteCloser) (*Session, error) {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.5"
	session.UserAgent = userAgent
	session.UserAgentPassword = userAgentPw
	session.HttpMethod = "GET"
	session.Accept = "*/*"

	transport := http.DefaultTransport
	if logger != nil {
		dial := WireLog(logger)
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: dial,

		}
	}
	retsTransport := RetsTransport{
		transport: transport,
		session: session,
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	session.Client = http.Client{
		Transport: &retsTransport,
		Jar: jar,
	}
	return &session, nil
}

/* wrapper to intercept each http call */
type RetsTransport struct {
	transport http.RoundTripper
	session Session
}

func (t *RetsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add(USER_AGENT, t.session.UserAgent)
	req.Header.Add(RETS_VERSION, t.session.Version)
	req.Header.Add(ACCEPT, t.session.Accept)

	if t.session.UserAgentPassword != "" {
		requestId := req.Header.Get(RETS_REQUEST_ID)
		sessionId := ""
		if h,err := req.Cookie(RETS_SESSION_ID); err == nil {
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
	challenge := res.Header.Get(WWW_AUTH)
	if !strings.HasPrefix(strings.ToLower(challenge), "digest") {
		return nil, errors.New("unknown authentication challenge: "+challenge)
	}
	req.Header.Add(WWW_AUTH_RESP, DigestResponse(challenge, t.session.Username, t.session.Password, req.Method, req.URL.Path))
	return t.transport.RoundTrip(req)
}



