/**
 */
package gorets

import (
	"net/http"
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

func NewSession(user, pw, userAgent, userAgentPw string) *Session {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.5"
	session.UserAgent = userAgent
	session.HttpMethod = "GET"
	session.Accept = "*/*"
	session.Client = &http.Client{}
	return &session
}

func (s *Session) createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(s.HttpMethod, url, nil)
	if err != nil {
		return nil, err
	}


	// TODO do all of this per request
	req.Header.Add(USER_AGENT, s.UserAgent)
	req.Header.Add(RETS_VERSION, s.Version)
	req.Header.Add(ACCEPT, s.Accept)

	return req, nil
}

type Session struct {
	Username,Password string
	UserAgentPassword string
	HttpMethod string
	Version string // "Threewide/1.5"

	Accept string
	UserAgent string

	Client *http.Client

}

type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

