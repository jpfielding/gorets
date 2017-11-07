package explorer

import (
	"context"
	"net/http"

	"github.com/jpfielding/gorets/config"
	"github.com/jpfielding/gorets/rets"
)

// ConnectionService ...
type ConnectionService struct {
	Connections map[string]config.Config
}

// ConnectionList ...
type ConnectionList struct {
	Connections []config.Config `json:"connections"`
}

// ConnectionListArgs ...
type ConnectionListArgs struct{}

// List ....
func (cs ConnectionService) List(r *http.Request, args *ConnectionListArgs, reply *ConnectionList) error {
	for _, v := range cs.Connections {
		reply.Connections = append(reply.Connections, v)
	}
	return nil
}

// AddConnectionArgs ..
type AddConnectionArgs struct {
	Connection config.Config `json:"connection"`
	Test       bool          `json:"test"`
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
		sess, err := args.Connection.Connect(ctx)
		err = sess.Process(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
			return nil
		})
		reply.Active = (err == nil)
		if err != nil {
			return err
		}
		reply.Tested = true
	}
	cs.Connections[args.Connection.ID] = args.Connection
	// TODO should we persist that change somewhere?
	reply.ID = args.Connection.ID
	return nil
}
