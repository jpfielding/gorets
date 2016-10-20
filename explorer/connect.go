package explorer

import (
	"context"
	"net/http"
)

// ConnectionService ...
type ConnectionService struct {
	connections map[string]Connection
}

// Load ...
func (cs *ConnectionService) Load() map[string]Connection {
	if cs.connections == nil {
		cs.connections = make(map[string]Connection)
		JSONLoad("/tmp/gorets/connections.json", &cs.connections)
	}
	return cs.connections
}

// Stash ..
func (cs *ConnectionService) Stash() {
	JSONStore("/tmp/gorets/connections.json", &cs.connections)
}

// ConnectionList ...
type ConnectionList struct {
	Connections []Connection
}

// ConnectionListArgs ...
type ConnectionListArgs struct {
	Active *bool `json:"active,omitempty"`
}

// List ....
func (cs ConnectionService) List(r *http.Request, args *ConnectionListArgs, reply *ConnectionList) error {
	for _, v := range cs.Load() {
		// if we want to filter on active
		if args.Active != nil {
			// do they have matching state
			if *args.Active != v.Active() {
				continue
			}
		}
		reply.Connections = append(reply.Connections, v)
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
	if args.Test {
		ctx := context.Background()
		if _, _, err := args.Connection.Login(ctx); err != nil {
			return err
		}
		reply.Tested = true
		reply.Active = args.Connection.Active()
	}
	reply.ID = args.Connection.ID
	cs.Load()
	cs.connections[args.Connection.ID] = args.Connection
	cs.Stash()
	return nil
}
