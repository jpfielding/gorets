package config

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"golang.org/x/net/proxy"
)

// Config ...
type Config struct {
	ID          string `json:"id"`
	LoginURL    string `json:"loginURL"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserAgent   string `json:"userAgent"`
	UserAgentPw string `json:"userAgentPw"`
	Proxy       string `json:"proxy"`
	RetsVersion string `json:"retsVersion"`
}

// Connect ...
// TODO need to remove connecting from session creation
func (c Config) Connect(ctx context.Context, wlog string) (*Session, error) {
	// start with the default Dialer from http.DefaultTransport
	transport := wirelog.NewHTTPTransport()
	// if there is a need to proxy
	if c.Proxy != "" {
		log.Printf("Using proxy %s", c.Proxy)
		d, err := proxy.SOCKS5("tcp", c.Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("rets proxy: %v", err)
		}
		transport.Dial = d.Dial
	}
	var closer io.Closer
	var err error
	// wire logging
	if wlog != "" {
		closer, err = wirelog.LogToFile(transport, wlog, true, true)
		if err != nil {
			return nil, fmt.Errorf("wirelog setup: %v", err)
		}
		log.Printf("wire logging enabled %s", wlog)
	}

	// should we throw an err here too?
	sess, err := rets.DefaultSession(
		c.Username,
		c.Password,
		c.UserAgent,
		c.UserAgentPw,
		c.RetsVersion,
		transport,
	)
	urls, err := rets.Login(ctx, sess, rets.LoginRequest{URL: c.LoginURL})
	if err != nil {
		return nil, err
	}
	return &Session{
		requester: sess,
		closer:    closer,
		urls:      *urls,
	}, err
}

// Session ...
type Session struct {
	Config Config // TODO should this be private?
	// things in user state
	closer    io.Closer
	requester rets.Requester
	urls      rets.CapabilityURLs
}

// Close is an io.Closer
// TODO remove need for context to help it match up with io.Closer
func (s *Session) Close(ctx context.Context) error {
	_, err := rets.Logout(ctx, s.requester, rets.LogoutRequest{URL: s.urls.Logout})
	if s.closer != nil {
		s.closer.Close()
	}
	s.closer = nil
	s.requester = nil
	return err
}

// Op is simple Rets function
type Op func(rets.Requester, rets.CapabilityURLs) error

// Process processes a set of requests
// TODO needs a retry mechanism
// TOOD deal with keeping sessions alive
func (s *Session) Process(ctx context.Context, ops ...Op) error {
	for _, op := range ops {
		err := op(s.requester, s.urls)
		if err != nil {
			return err
		}
	}
	return nil
}
