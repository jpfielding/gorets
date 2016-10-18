package server

import (
	"context"
	"encoding/json"
	"fmt"
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
func Metadata(ctx context.Context, u User) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Request == nil {
			http.Error(w, "Not Logged in", 400)
			return
		}
		compact := &retsutil.IncrementalCompact{}
		err = compact.Load(s.Requester, ctx, u.URLs.GetMetadata)
		if err != nil {
			fmt.Println("extracting metadata", err)
		}
		standard := retsutil.AsStandard(*compact).Convert()
		json.NewEncoder(w).Encode(standard)
	}
}
