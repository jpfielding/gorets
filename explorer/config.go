package explorer

import (
	"context"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// ConfigService ...
type ConfigService struct {
	Configs map[string]Config
}

// ConfigList ...
type ConfigList struct {
	Configs []Config `json:"configs"`
}

// ConfigListArgs ...
type ConfigListArgs struct{}

// List ....
func (cs ConfigService) List(r *http.Request, args *ConfigListArgs, reply *ConfigList) error {
	for _, v := range cs.Configs {
		reply.Configs = append(reply.Configs, v)
	}
	return nil
}

// AddConfigArgs ..
type AddConfigArgs struct {
	Config Config `json:"config"`
	Test   bool   `json:"test"`
}

// AddConfigReply ...
type AddConfigReply struct {
	ID     string `json:"id"`
	Tested bool   `json:"tested"`
	Active bool   `json:"active"`
}

// Add ....
func (cs ConfigService) Add(r *http.Request, args *AddConfigArgs, reply *AddConfigReply) error {
	if args.Test {
		ctx := context.Background()
		sess, err := args.Config.Connect(ctx)
		err = sess.Process(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
			return nil
		})
		reply.Active = (err == nil)
		if err != nil {
			return err
		}
		reply.Tested = true
	}
	cs.Configs[args.Config.ID] = args.Config
	// TODO should we persist that change somewhere?
	reply.ID = args.Config.ID
	return nil
}
