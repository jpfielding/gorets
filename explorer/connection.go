package explorer

import (
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/config"
)

// TODO remove this connection service and replace with configservice (or reference the service from the react)

// ConnectionService ...
type ConnectionService struct {
	Client config.Client
}

// List ....
func (cs *ConnectionService) List(r *http.Request, args *config.ListArgs, reply *config.ListReply) error {
	tmp, err := cs.Client.List(*args)
	reply.Configs = tmp.Configs
	fmt.Printf("connection.List found %d connections\n", len(reply.Configs))
	return err
}
