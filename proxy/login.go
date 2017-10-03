package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"golang.org/x/net/proxy"
)

// Login manages de/multiplexing requests to RETS servers
func Login(prefix string, proxy []Config) http.HandlerFunc {
	srcs := sources{}
	for _, c := range proxy {
		if _, ok := srcs[c.Service]; !ok {
			srcs[c.Service] = logins{}
		}
		srcs[c.Service][c.User] = login{config: c}
	}
	return func(res http.ResponseWriter, req *http.Request) {
		sub := strings.TrimPrefix(req.URL.Path, prefix)
		parts := strings.Split(sub, "/")
		src := parts[0]
		usr := parts[1]
		if _, ok := srcs[src]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "source %s not found", src)
			return
		}
		if _, ok := srcs[src][usr]; !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "user %s not found", usr)
			return
		}
		fmt.Fprintf(res, "Hi there, I love proxying %s for %s login!", src, usr)
	}
}

type sources map[string]logins

type logins map[string]login

// login holds the state for an established session
type login struct {
	config  Config
	session rets.Requester
	urls    *rets.CapabilityURLs
	closer  io.Closer // holds the wirelog output
}

func (l *login) Clear() {
	if l.session == nil {
		return
	}
	ctx := context.Background()
	req := rets.LogoutRequest{URL: l.urls.Logout}
	rets.Logout(ctx, l.session, req)
	if l.closer != nil {
		l.closer.Close()
	}
	l.session = nil
}

// Session returns the cached rets session
func (l *login) Session() (rets.Requester, *rets.CapabilityURLs, error) {
	if l.session == nil {
		session, closer, urls, err := l.createSession()
		if err != nil {
			return nil, nil, fmt.Errorf("rets session create")
		}
		l.session = session
		l.urls = urls
		l.closer = closer
	}
	return l.session, l.urls, nil
}
func (l *login) createSession() (rets.Requester, io.Closer, *rets.CapabilityURLs, error) {
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

// Config sets up a conneciton
// TODO does any auth require the version info for auth
type Config struct {
	Service                      string
	URL                          string
	User, Password               string
	UserAgent, UserAgentPassword string
	Version                      string
	Proxy                        string
}
