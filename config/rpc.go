package config

import (
	"net/http"
)

// ConfigService ...
type ConfigService struct {
	Configs map[string]Config
}

// ListReply ...
type ListReply struct {
	Configs []Config `json:"configs"`
}

// CListArgs ...
type ListArgs struct{}

// List ....
// curl http://localhost:8888/rpc -H "Content-Type: application/json" -X POST -d '{"id": 1, "method":"ConfigService.List","params":[{}]}'
func (cs *ConfigService) List(r *http.Request, args *ListArgs, reply *ListReply) error {
	for _, v := range cs.Configs {
		reply.Configs = append(reply.Configs, v)
	}
	return nil
}

// ConfigAddArgs ..
type AddArgs struct {
	Config Config `json:"Config"`
}

// ConfigAddReply ...
type AddReply struct {
	Config Config `json:"Config"`
}

// Add ....
func (cs *ConfigService) Add(r *http.Request, args *AddArgs, reply *AddReply) error {
	cs.Configs[args.Config.ID] = args.Config
	// TODO should we persist that change somewhere?
	reply.Config = args.Config
	return nil
}
