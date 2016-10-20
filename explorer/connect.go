package explorer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ConnectionService ...
type ConnectionService struct {
	connections map[string]Connection
}

// Load ...
func (cs *ConnectionService) Load() map[string]Connection {
	if len(cs.connections) == 0 {
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

// List ....
func (cs ConnectionService) List(r *http.Request, args *struct{}, reply *ConnectionList) error {
	for _, v := range cs.Load() {
		reply.Connections = append(reply.Connections, v)
	}
	return nil
}

// AddConnectionArgs ..
type AddConnectionArgs struct {
	Connection Connection
	Test       bool
}

// Add ....
func (cs ConnectionService) Add(r *http.Request, args *AddConnectionArgs, reply *struct{}) error {
	if args.Test {
		ctx := context.Background()
		if _, err := args.Connection.Login(ctx); err != nil {
			return err
		}
	}
	cs.Load()
	cs.connections[args.Connection.ID] = args.Connection
	cs.Stash()
	return nil
}

// Connect ...
// input: Connection
// output: rets.CapabilityURLS
func Connect() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Connection
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Printf("params: %v\n", p)
		ctx := context.Background()
		_, err = p.Login(ctx)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		cs := ConnectionService{}
		cs.Load()
		cs.connections[p.ID] = p
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p.URLs)

		cs.Stash()

	}
}
