/**
 */
package gorets

import (
	"errors"
	"net/http"
	"strings"
)

const DEFAULT_TIMEOUT int = 300000

/** header values */
const (
	RETS_VERSION string = "RETS-Version"
	RETS_SESSION_ID string = "RETS-Session-ID"
	RETS_REQUEST_ID string = "RETS-Request-ID"
	USER_AGENT string = "User-Agent"
	RETS_UA_AUTH_HEADER string = "RETS-UA-Authorization"
	ACCEPT string = "Accept"
	ACCEPT_ENCODING string = "Accept-Encoding"
	CONTENT_ENCODING string = "Content-Encoding"
	DEFLATE_ENCODINGS string = "gzip,deflate"
	CONTENT_TYPE string = "Content-Type"
	WWW_AUTH string = "Www-Authenticate"
	WWW_AUTH_RESP string = "Authorization"

)

type Session struct {
	Username,Password string

	UserAgent, UserAgentPassword string
	HttpMethod string

	Version string // "Threewide/1.5"

	Accept string

	Client *http.Client

}

type RetsTransport struct {
	transport http.RoundTripper
	session *Session
}

type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

func NewSession(user, pw, userAgent, userAgentPw string) *Session {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.5"
	session.UserAgent = userAgent
	session.HttpMethod = "GET"
	session.Accept = "*/*"


	session.Client = &http.Client{
		Transport: &RetsTransport{session: &session},
	}
	return &session
}

func (t *RetsTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add(USER_AGENT, t.session.UserAgent)
	req.Header.Add(RETS_VERSION, t.session.Version)
	req.Header.Add(ACCEPT, t.session.Accept)

	if t.session.UserAgentPassword != "" {
		requestId := resp.Header.Get(RETS_REQUEST_ID)
		sessionId,err := req.Cookie(RETS_SESSION_ID)
		if err != nil {
			return nil, err
		}
		uaAuthHeader := CalculateUaAuthHeader(
			t.session.UserAgent,
			t.session.UserAgentPassword,
			requestId,
			sessionId.Value,
			t.session.Version,
		)
		req.Header.Add(RETS_UA_AUTH_HEADER, uaAuthHeader)
	}

	res, err := t.transport.RoundTrip(req)
	if res.StatusCode == 401 {
		challenge := resp.Header.Get(WWW_AUTH)
		if !strings.HasPrefix(strings.ToLower(challenge), "digest") {
			return nil, errors.New("unknown authentication challenge: "+challenge)
		}
		req.Header.Add(WWW_AUTH_RESP, DigestResponse(challenge, t.session.Username, t.session.Password, req.Method, req.URL.Path))
	}

	return t.transport.RoundTrip(req)
}



