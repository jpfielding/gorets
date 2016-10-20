package explorer

import (
	"context"
	"fmt"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// TODO make this a user based map of sessions
var sessions = Sessions{}

// Sessions ...
type Sessions map[string]Session

// Open ...
func (s Sessions) Open(id string) Session {
	if _, ok := s[id]; !ok {
		s[id] = Session{Connection: connections[id]}
	}
	return s[id]
}

// Session ...
type Session struct {
	Connection Connection
	// Requester is user state
	requester rets.Requester
	// URLs need this to know where to route requests
	urls *rets.CapabilityURLs
}

// Wirelog path
func (c *Session) Wirelog() string {
	return fmt.Sprintf("/tmp/gorets/%s/wire.log", c.Connection.ID)
}

// MSystem path
func (c *Session) MSystem() string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.Connection.ID)
}

// session ...
func (c *Session) session() (rets.Requester, error) {
	if c.requester != nil {
		return c.requester, nil
	}
	// start with the default Dialer from http.DefaultTransport
	transport := wirelog.NewHTTPTransport()
	// logging
	err := wirelog.LogToFile(transport, c.Wirelog(), true, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("wire logging enabled:", c.Wirelog())
	conn := c.Connection
	r, err := rets.DefaultSession(
		conn.Username,
		conn.Password,
		conn.UserAgent,
		conn.UserAgentPw,
		conn.Version,
		transport,
	)
	c.requester = r
	return r, err
}

// Active tells whether this connection is considered to be in use
func (c *Session) Active() bool {
	return c.urls != nil
}

//Login ...
func (c *Session) Login(ctx context.Context) (rets.Requester, *rets.CapabilityURLs, error) {
	r, err := c.session()
	if err != nil {
		return nil, nil, err
	}
	if c.urls != nil {
		return r, c.urls, nil
	}
	fmt.Printf("login: %v\n", c.Connection.URL)
	urls, err := rets.Login(r, ctx, rets.LoginRequest{URL: c.Connection.URL})
	c.urls = urls
	return r, urls, err

}
