package config

import (
	"net/http"
)

// RPCService ...
type RPCService struct {
	Configs func(*ListArgs) ([]Config, error)
}

// ListReply ...
type ListReply struct {
	Configs []Config `json:"configs"`
}

// ListArgs ...
type ListArgs struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Active *bool  `json:"active"`
}

// List ....
// curl http://localhost:8888/rpc -H "Content-Type: application/json" -X POST -d '{"id": 1, "method":"ConfigService.List","params":[{}]}'
func (cs *RPCService) List(r *http.Request, args *ListArgs, reply *ListReply) error {
	cfgs, err := cs.Configs(args)
	if err != nil {
		return err
	}
	for _, v := range cfgs {
		reply.Configs = append(reply.Configs, v)
	}
	return nil
}

// TODO add a save func
