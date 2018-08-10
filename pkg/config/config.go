package config

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jpfielding/gorets/pkg/rets"
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
func (c Config) Connect(ctx context.Context, wlog io.Writer) (*Session, error) {
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
	// wire logging
	if wlog != nil {
		err := wirelog.LogToWriter(transport, wlog, true, true)
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
	if err != nil {
		return nil, err
	}
	urls, err := rets.Login(ctx, sess, rets.LoginRequest{URL: c.LoginURL})
	if err != nil {
		return nil, err
	}
	closer := func() error {
		_, err := rets.Logout(ctx, sess, rets.LogoutRequest{URL: urls.Logout})
		return err
	}
	return &Session{
		requester: sess,
		close:     closer,
		urls:      *urls,
	}, err
}

// Session ...
type Session struct {
	Config Config // TODO should this be private?
	// things in user state
	close     func() error
	requester rets.Requester
	urls      rets.CapabilityURLs
}

// Close is an io.Closer
func (s *Session) Close() error {
	var err error
	if s.close != nil {
		err = s.close()
	}
	s.close = nil
	s.requester = nil
	return err
}

// Op is simple Rets function
type Op func(rets.Requester, rets.CapabilityURLs) error

// Process processes a set of requests
func (s *Session) Process(ctx context.Context, ops ...Op) error {
	for _, op := range ops {
		//op = retry(op, 3)
		err := op(s.requester, s.urls)
		if err != nil {
			return err
		}
	}
	return nil
}

// retry operation
func retry(op Op, times int) Op {
	return func(req rets.Requester, urls rets.CapabilityURLs) error {
		var err error
		for ; times > 0; times-- {
			err = op(req, urls)
			if err == nil {
				return nil
			}
		}
		return err
	}
}
