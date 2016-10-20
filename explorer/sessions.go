package explorer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// TODO make this a user based map of sessions
var sessions = Sessions{}

// Sessions ...
type Sessions map[string]*Session

// Open ...
func (s Sessions) Open(id string) *Session {
	if _, ok := s[id]; !ok {
		s[id] = &Session{Connection: connections[id]}
	}
	return s[id]
}

// Session ...
type Session struct {
	Connection Connection
	// Requester is user state
	transport *http.Transport
	requester rets.Requester
	urls      rets.CapabilityURLs
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
func (c *Session) create() (rets.Requester, error) {
	if c.transport == nil {
		// start with the default Dialer from http.DefaultTransport
		transport := wirelog.NewHTTPTransport()
		// logging
		err := wirelog.LogToFile(transport, c.Wirelog(), true, true)
		if err != nil {
			return nil, err
		}
		fmt.Println("wire logging enabled:", c.Wirelog())
		c.transport = transport
	}
	conn := c.Connection
	r, err := rets.DefaultSession(
		conn.Username,
		conn.Password,
		conn.UserAgent,
		conn.UserAgentPw,
		conn.Version,
		c.transport,
	)
	return r, err
}

// Active tells whether this connection is considered to be in use
func (c *Session) Active() bool {
	return c.transport != nil
}

// Close is a closer
func (c *Session) Close() error {
	// TODO see if we can close the file handle for the wirelog
	c.transport = nil
	return nil
}

// Exec ...
func (c *Session) Exec(ctx context.Context, op func(rets.Requester, rets.CapabilityURLs) error) error {
	if c.requester == nil {
		r, err := c.create()
		urls, err := rets.Login(r, ctx, rets.LoginRequest{URL: c.Connection.URL})
		if err != nil {
			return err
		}
		c.requester = r
		c.urls = *urls
		// defer rets.Logout(r, ctx, rets.LogoutRequest{URL: urls.Logout})
	}
	return op(c.requester, c.urls)
}
