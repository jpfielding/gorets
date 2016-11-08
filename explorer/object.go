package explorer

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// ObjectParams ...
type ObjectParams struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Type     string `json:"type"`
	ObjectID string `json:"objectid"`
	Location int    `json:"location"` // setting to 1 requests the URL to the photo
}

// Objects ...
type Objects struct {
	Objects []Object
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
	log.Printf("object get params: %v\n", args)

	s := sessions.Open(args.ID)
	if s == nil {
		return fmt.Errorf("no source found for %s", args.ID)
	}
	ctx := context.Background()
	return s.Exec(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		// warning, this does _all_ of the photos
		response, err := rets.GetObjects(r, ctx, rets.GetObjectRequest{
			URL: u.GetObject,
			GetObjectParams: rets.GetObjectParams{
				Resource: args.Resource,
				Type:     args.Type,
				ID:       args.ObjectID,
			},
		})
		if err != nil {
			return err
		}
		// open the json encoder
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
			return nil
		})
	})
}
