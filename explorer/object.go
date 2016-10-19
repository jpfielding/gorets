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
	Blob       string            `json:",omitempty"`
}

// GetObject ...
// input: ObjectParams
// output: map[string]Object <- streaming
// {
// 	"2927498:1": {"ContentID":"2927498","ContentType":"image/jpeg","ObjectID":1},
// 	"2927498:2": {"ContentID":"2927498","ContentType":"image/jpeg","ObjectID":2},
// 	"2927498:3": {"ContentID":"2927498","ContentType":"image/jpeg","ObjectID":3}
// }
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
		w.Write([]byte("{"))
		// open the json encoder
		enc := json.NewEncoder(w)
		defer response.Close()
		comma := false
		response.ForEach(func(o *rets.Object, err error) error {
			if comma {
				w.Write([]byte(","))
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
				Blob:           base64.StdEncoding.EncodeToString(o.Blob),
			}
			w.Write([]byte(fmt.Sprintf("\"%s:%d\": ", o.ContentID, o.ObjectID)))
			enc.Encode(obj)
			comma = true
			return err
		})
		w.Write([]byte("}"))
	}
}
