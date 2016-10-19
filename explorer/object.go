package explorer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jpfielding/gorets/rets"
)

// ObjectParams ...
type ObjectParams struct {
	Resource string `json:"resource"`
	Type     string `json:"type"`
	ID       string `json:"id"`
	Location int    `json:"location"` // setting to 1 requests the URL to the photo
}

// Object ...
type Object struct {
	ContentID,
	ContentType string
	ObjectID int
	UID      string
	Description,
	SubDescription,
	Location string
	RetsError bool
	// RetsMessage *rets.Response
	Preferred  bool
	ObjectData map[string]string
	Blob       string
}

// GetObject ...
// input: ObjectParams
// output: []Objects
func GetObject(ctx context.Context, c *Connection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p ObjectParams
		if r.Body == nil {
			http.Error(w, "missing object params", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Printf("params: %v\n", p)

		// warning, this does _all_ of the photos
		response, err := rets.GetObjects(c.Requester, ctx, rets.GetObjectRequest{
			URL:      c.URLs.GetObject,
			Resource: p.Resource,
			Type:     p.Type,
			ID:       p.ID,
		})
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		defer response.Close()
		response.ForEach(func(o *rets.Object, err error) error {
			obj := Object{
				ContentID:      o.ContentID,
				ContentType:    o.ContentType,
				UID:            o.UID,
				Description:    o.Description,
				SubDescription: o.SubDescription,
				Location:       o.Location,
				RetsError:      o.RetsError,
				Preferred:      o.Preferred,
				ObjectData:     o.ObjectData,
				Blob:           base64.StdEncoding.EncodeToString(o.Blob),
			}
			enc.Encode(obj)
			return err
		})
	}
}
