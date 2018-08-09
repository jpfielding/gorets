package proxy

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"golang.org/x/net/proxy"
)

// Config sets up a connection
type Config struct {
	Service                      string
	URL                          string
	User, Password               string
	UserAgent, UserAgentPassword string
	Version                      string
	Proxy                        string
}

// Sources is a indexer for looking up sessions
type Sources map[string]Sessions

// NewSources creates the index for sessions/users
func NewSources(cfgs []Config) Sources {
	srcs := Sources{}
	for _, c := range cfgs {
		if _, ok := srcs[c.Service]; !ok {
			srcs[c.Service] = Sessions{}
		}
		srcs[c.Service][c.User] = Session{config: c}
	}
	return srcs
}

// Sessions a wrapper to manage requests
type Sessions map[string]Session

// Session holds the state for an established session
type Session struct {
	config    Config
	requester rets.Requester
	urls      *rets.CapabilityURLs
	closer    io.Closer // holds the wirelog output
}

// Clear the current session
func (l *Session) Clear() {
	if l.requester == nil {
		return
	}
	ctx := context.Background()
	req := rets.LogoutRequest{URL: l.urls.Logout}
	rets.Logout(ctx, l.requester, req)
	if l.closer != nil {
		l.closer.Close()
	}
	l.requester = nil
}

// Get returns the cached rets session
func (l *Session) Get() (rets.Requester, *rets.CapabilityURLs, error) {
	if l.requester == nil {
		req, closer, urls, err := l.create()
		if err != nil {
			return nil, nil, fmt.Errorf("rets session create")
		}
		l.requester = req
		l.urls = urls
		l.closer = closer
	}
	return l.requester, l.urls, nil
}

func (l *Session) create() (rets.Requester, io.Closer, *rets.CapabilityURLs, error) {
	// should we throw an err here too?
	session, closer, err := NewSession(l.config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("rets session setup")
	}
	// login and get our urls
	req := rets.LoginRequest{URL: l.config.URL}
	ctx := context.Background()
	capability, err := rets.Login(ctx, session, req)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("rets session login")
	}
	return session, closer, capability, nil
}

// NewSession creates a Rets session from the given config
func NewSession(c Config) (rets.Requester, io.Closer, error) {
	// start with the default Dialer from http.DefaultTransport
	transport := wirelog.NewHTTPTransport()
	// if there is a need to proxy
	if c.Proxy != "" {
		log.Printf("Using proxy %s", c.Proxy)
		d, err := proxy.SOCKS5("tcp", c.Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, nil, fmt.Errorf("rets proxy: '%s'", c.Proxy)
		}
		transport.Dial = d.Dial
	}
	var closer io.Closer
	var err error
	// wire logging
	logFile := fmt.Sprintf("/tmp/rets/wirelog/%s-%s.log", c.Service, c.User)
	closer, err = wirelog.LogToFile(transport, logFile, true, true)
	if err != nil {
		return nil, nil, fmt.Errorf("wirelog setup")
	}
	log.Printf("wire logging enabled %s", logFile)

	// should we throw an err here too?
	sess, err := rets.DefaultSession(
		c.User,
		c.Password,
		c.UserAgent,
		c.UserAgentPassword,
		c.Version,
		transport,
	)
	return sess, closer, err
}
