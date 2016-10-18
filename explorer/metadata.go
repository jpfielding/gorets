package explorer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gorets/retsutil"
)

// MetadataResponse ...
type MetadataResponse struct {
	Metadata metadata.MSystem
}

// MetadataRequestType is a typedef metadata extraction options
type MetadataRequestType func(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error)

// options for extracting metadata
var options = map[string]MetadataRequestType{
	"STANDARD-XML":        getStandardMetadata,
	"COMPACT":             getCompactMetadata,
	"COMPACT-INCREMENTAL": getCompactIncremental,
}

// Metadata ...
// input: query param extraction=(|STANDARD-XML|COMPACT|COMPACT-INCREMENTAL)
// output: metadata.MSystem
func Metadata(ctx context.Context, c *Connection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.Requester == nil {
			http.Error(w, "Not Logged in", 400)
			return
		}
		extraction := r.URL.Query().Get("extraction")
		if extraction == "" {
			extraction = "COMPACT"
		}
		if op, ok := options[extraction]; ok {
			// lookup the operation for pulling metadata
			standard, err := op(c.Requester, ctx, c.URLs.GetMetadata)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(standard)
		} else {
			http.Error(w, fmt.Sprintf("%s not supported", extraction), 400)
			return
		}
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
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	compact, err := rets.ParseMetadataCompactResult(reader)
	if err != nil {
		return nil, err
	}
	return retsutil.AsStandard(*compact).Convert()
}

// getStandardMetadata ...
func getStandardMetadata(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
	reader, err := rets.MetadataStream(requester, ctx, rets.MetadataRequest{
		URL:    url,
		Format: "STANDARD-XML",
		MType:  "METADATA-SYSTEM",
		ID:     "*",
	})
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	parser := xml.NewDecoder(reader)
	rets := metadata.RETSResponseWrapper{}
	err = parser.Decode(&rets)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return &rets.Metadata.MSystem, err
}
