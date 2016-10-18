package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/retsutil"
)

// MetadataResponse ...
type MetadataResponse struct {
	Metadata metadata.MSystem
}

// Metadata ...
// input:
// output: metadata.MSystem
func Metadata(ctx context.Context, u *User) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if u.Requester == nil {
			http.Error(w, "Not Logged in", 400)
			return
		}
		compact := &retsutil.IncrementalCompact{}
		err := compact.Load(u.Requester, ctx, u.URLs.GetMetadata)
		if err != nil {
			http.Error(w, "metadata request failed", 400)
			return
		}
		standard, err := retsutil.AsStandard(*compact).Convert()
		if err != nil {
			http.Error(w, "metadata conversion failed", 400)
			return
		}
		json.NewEncoder(w).Encode(standard)
	}
}
