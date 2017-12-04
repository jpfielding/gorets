package explorer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/config"
	"github.com/jpfielding/gorets/rets"
)

// ObjectParams ...
type ObjectParams struct {
	Connection config.Config `json:"connection"`
	Resource   string        `json:"resource"`
	Type       string        `json:"type"`
	ID         string        `json:"ids"`
	Location   int           `json:"location"` // setting to 1 requests the URL to the photo
}

// Objects ...
type Objects struct {
	Objects []Object `json:"Objects"`
	Wirelog string   `json:"wirelog,omitempty"`
}

// Object ...
type Object struct {
	ContentID      string `json:",omitempty"`
	ContentType    string `json:",omitempty"`
	ObjectID       int    `json:",omitempty"`
	UID            string `json:",omitempty"`
	Description    string `json:",omitempty"`
	SubDescription string `json:",omitempty"`
	Location       string `json:",omitempty"`
	RetsError      bool   `json:",omitempty"`
	// RetsMessage *rets.Response
	Preferred  bool              `json:",omitempty"`
	ObjectData map[string]string `json:",omitempty"`
	Blob       []byte            `json:",omitempty"`
}

// ObjectService ...
type ObjectService struct{}

// Get ....
func (os ObjectService) Get(r *http.Request, args *ObjectParams, reply *Objects) error {
	fmt.Printf("object get params: %v\n", args)

	cfg := args.Connection
	ctx := context.Background()
	wirelog := bytes.Buffer{}
	sess, err := cfg.Connect(ctx, &wirelog)
	if err != nil {
		return err
	}
	defer sess.Close()
	return sess.Process(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		// warning, this does _all_ of the photos
		rsp, err := rets.GetObjects(ctx, r, rets.GetObjectRequest{
			URL: u.GetObject,
			GetObjectParams: rets.GetObjectParams{
				Resource: args.Resource,
				Type:     args.Type,
				ID:       args.ID,
			},
		})
		if err != nil {
			reply.Wirelog = string(wirelog.Bytes())
			return err
		}
		// open the json encoder
		response := &rets.GetObjectResponse{Response: rsp}
		defer response.Close()
		return response.ForEach(func(o *rets.Object, err error) error {
			if o == nil {
				return err
			}
			// translate
			obj := Object{
				ContentID:      o.ContentID,
				ContentType:    o.ContentType,
				ObjectID:       o.ObjectID,
				UID:            o.UID,
				Description:    o.Description,
				SubDescription: o.SubDescription,
				Location:       o.Location,
				RetsError:      o.RetsError,
				Preferred:      o.Preferred,
				ObjectData:     o.ObjectData,
				Blob:           o.Blob,
			}
			reply.Objects = append(reply.Objects, obj)
			reply.Wirelog = string(wirelog.Bytes())
			return nil
		})
	})
}
