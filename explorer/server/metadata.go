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
func Metadata(s Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		compact := &retsutil.IncrementalCompact{}
		// TODO deal with contexts in the web appropriately
		ctx := context.Background()
		err = compact.Load(s.Requester, ctx, s.URLs.GetMetadata)
		if err != nil {
			fmt.Println("extracting metadata", err)
		}
		standard := retsutil.AsStandard(*compact).Convert()
		json.NewEncoder(w).Encode(standard)
	}
}
