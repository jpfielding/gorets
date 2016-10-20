package explorer

import (
	"context"
	"fmt"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// Connection ...
type Connection struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserAgent   string `json:"user-agent"`
	UserAgentPw string `json:"user-agent-pw"`
	Version     string `json:"rets-version"`

	// Requester is user state
	requester rets.Requester
	// URLs need this to know where to route requests
	urls *rets.CapabilityURLs
}

// Wirelog path
func (c *Connection) Wirelog() string {
	return fmt.Sprintf("/tmp/gorets/%s/wire.log", c.ID)
}

// MSystem path
func (c *Connection) MSystem() string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.ID)
}

// session ...
func (c *Connection) session() (rets.Requester, error) {
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
	r, err := rets.DefaultSession(
		c.Username,
		c.Password,
		c.UserAgent,
		c.UserAgentPw,
		c.Version,
		transport,
	)
	c.requester = r
	return r, err
}

// Active tells whether this connection is considered to be in use
func (c *Connection) Active() bool {
	return c.urls != nil
}

//Login ...
func (c *Connection) Login(ctx context.Context) (rets.Requester, *rets.CapabilityURLs, error) {
	r, err := c.session()
	if err != nil {
		return nil, nil, err
	}
	if c.urls != nil {
		return r, c.urls, nil
	}
	fmt.Printf("login: %v\n", c.URL)
	urls, err := rets.Login(r, ctx, rets.LoginRequest{URL: c.URL})
	c.urls = urls
	return r, urls, err

}
