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
	URLs *rets.CapabilityURLs `json:"-"`
}

// Wirelog path
func (c *Connection) Wirelog() string {
	return fmt.Sprintf("/tmp/gorets/%s/wire.log", c.ID)
}

// MSystem path
func (c *Connection) MSystem() string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.ID)
}

// Requester ...
func (c *Connection) Requester() (rets.Requester, error) {
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

//Login ...
func (c *Connection) Login(ctx context.Context) (rets.Requester, error) {
	r, err := c.Requester()
	if err != nil {
		return nil, err
	}
	if c.URLs != nil {
		return r, nil
	}
	fmt.Printf("login: %v\n", c.URL)
	urls, err := rets.Login(r, ctx, rets.LoginRequest{URL: c.URL})
	c.URLs = urls
	return r, err

}
