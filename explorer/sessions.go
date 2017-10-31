package explorer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"golang.org/x/net/proxy"
)

// TODO make this a user based map of sessions
var sessions = Sessions{}

// Sessions ...
type Sessions map[string]*Session

// Open ...
func (s Sessions) Open(id string) *Session {
	if _, ok := s[id]; !ok {
		if _, ok = connections[id]; !ok {
			return nil
		}
		s[id] = &Session{Connection: connections[id]}
	}
	return s[id]
}

// Session ...
type Session struct {
	Connection Connection
	// things in user state
	closer    io.Closer
	requester rets.Requester
	urls      rets.CapabilityURLs
}

// ReadWirelog ...
func (c *Session) ReadWirelog(fun func(*os.File, error) error) error {
	f, err := os.Open(c.Wirelog())
	defer f.Close()
	if e := fun(f, err); e != nil {
		return e
	}
	return nil
}

// Wirelog path
func (c *Session) Wirelog() string {
	return fmt.Sprintf("/tmp/gorets/%s/wire.log", c.Connection.ID)
}

// MSystem path
func (c *Session) MSystem() string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.Connection.ID)
}

// create ...
func (c *Session) create() (rets.Requester, io.Closer, error) {
	// start with the default Dialer from http.DefaultTransport
	transport := wirelog.NewHTTPTransport()
	// if there is a need to proxy
	if c.Connection.Proxy != "" {
		log.Printf("Using proxy %s", c.Connection.Proxy)
		d, err := proxy.SOCKS5("tcp", c.Connection.Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, nil, fmt.Errorf("rets proxy: '%s'", c.Connection.Proxy)
		}
		transport.Dial = d.Dial
	}
	var closer io.Closer
	var err error
	// wire logging
	logFile := c.Wirelog()
	closer, err = wirelog.LogToFile(transport, logFile, true, true)
	if err != nil {
		return nil, nil, fmt.Errorf("wirelog setup")
	}
	log.Printf("wire logging enabled %s", logFile)

	// should we throw an err here too?
	sess, err := rets.DefaultSession(
		c.Connection.Username,
		c.Connection.Password,
		c.Connection.UserAgent,
		c.Connection.UserAgentPw,
		c.Connection.Version,
		transport,
	)
	return sess, closer, err
}

// Active tells whether this connection is considered to be in use
func (c *Session) Active() bool {
	return c.requester != nil
}

// Close is a closer
func (c *Session) Close() error {
	defer c.closer.Close()
	// TODO see if we can close the file handle for the wirelog
	c.closer = nil
	c.requester = nil
	return nil
}

// Exec TODO needs a retry mechanism
func (c *Session) Exec(ctx context.Context, op func(rets.Requester, rets.CapabilityURLs) error) error {
	if c.requester == nil {
		r, closer, err := c.create()
		urls, err := rets.Login(ctx, r, rets.LoginRequest{URL: c.Connection.URL})
		if err != nil {
			return err
		}
		c.closer = closer
		c.requester = r
		c.urls = *urls
	}
	return op(c.requester, c.urls)
}
