package explorer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// TODO review GLOBAL state even though its a 'resource'
var connections = Connections{}.Load()

// Connections ...
type Connections map[string]Connection

// Load ...
func (cs Connections) Load() Connections {
	if len(cs) == 0 {
		fmt.Println("loading connections")
		JSONLoad("/tmp/gorets/connections.json", &cs)
	}
	fmt.Printf("found %d connections\n", len(cs))
	return cs
}

// Stash ..
func (cs Connections) Stash() {
	JSONStore("/tmp/gorets/connections.json", cs)
}

// Connection ...
type Connection struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserAgent   string `json:"user-agent"`
	UserAgentPw string `json:"user-agent-pw"`
	Version     string `json:"rets-version"`
}

// ConnectionService ...
type ConnectionService struct{}

// ConnectionList ...
type ConnectionList struct {
	Connections []Connection `json:"connections"`
}

// ConnectionListArgs ...
type ConnectionListArgs struct {
	Active *bool `json:"active,omitempty"`
}

// List ....
func (cs ConnectionService) List(r *http.Request, args *ConnectionListArgs, reply *ConnectionList) error {
	for _, v := range connections {
		// if we want to filter on active
		if args.Active != nil {
			// do they have matching state
			if _, ok := sessions[v.ID]; *args.Active != ok {
				continue
			}
		}
		reply.Connections = append(reply.Connections, v)
	}
	return nil
}

// DeleteConnectionArgs ..
type DeleteConnectionArgs struct {
	ID     string `json:"id"`
	Logout bool   `json:"logout"`
}

// Delete ...
func (cs ConnectionService) Delete(r *http.Request, args *DeleteConnectionArgs, reply *struct{}) error {
	delete(connections, args.ID)
	connections.Stash()
	if session, ok := sessions[args.ID]; !ok {
		if session.Active() && args.Logout {
			ctx := context.Background()
			sess, urls, err := session.Login(ctx)
			if err != nil {
				rets.Logout(sess, ctx, rets.LogoutRequest{URL: urls.Logout})
			}
			delete(sessions, args.ID)
		}
		return nil
	}

	return nil
}

// AddConnectionArgs ..
type AddConnectionArgs struct {
	Connection Connection `json:"connection"`
	Test       bool       `json:"test"`
}

// AddConnectionReply ...
type AddConnectionReply struct {
	ID     string `json:"id"`
	Tested bool   `json:"tested"`
	Active bool   `json:"active"`
}

// Add ....
func (cs ConnectionService) Add(r *http.Request, args *AddConnectionArgs, reply *AddConnectionReply) error {
	session := Session{Connection: args.Connection}
	if args.Test {
		ctx := context.Background()
		if _, _, err := session.Login(ctx); err != nil {
			return err
		}
		reply.Tested = true
	}
	connections[args.Connection.ID] = args.Connection
	connections.Stash()
	reply.Active = session.Active()
	reply.ID = args.Connection.ID
	return nil
}
