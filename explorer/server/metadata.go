package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
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
		var standard *metadata.MSystem
		var err error
		if _, ok := r.URL.Query()["incremental"]; ok {
			standard, err = getCompactIncremental(u.Requester, ctx, u.URLs.GetMetadata)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		standard, err = getCompactMetadata(u.Requester, ctx, u.URLs.GetMetadata)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(standard)
	}
}

// TODO extract a common func and switch on the incoming param

// getCompactIncremental retrieve the RETS Compact metadata from the server
func getCompactIncremental(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
	compact := &retsutil.IncrementalCompact{}
	err := compact.Load(requester, ctx, url)
	if err != nil {
		return nil, err
	}
	return retsutil.AsStandard(*compact).Convert()
}

// getCompactMetadata retrieve the RETS Compact metadata from the server
func getCompactMetadata(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
	reader, err := rets.MetadataStream(requester, ctx, rets.MetadataRequest{
		URL:    url,
		Format: "COMPACT",
		MType:  "METADATA-SYSTEM",
		ID:     "*",
	})
	if err != nil {
		return nil, err
	}
	compact, err := rets.ParseMetadataCompactResult(reader)
	if err != nil {
		return nil, err
	}
	return retsutil.AsStandard(*compact).Convert()
}
