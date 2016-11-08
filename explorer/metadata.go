package explorer

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gorets/util"
)

// MetadataService ...
type MetadataService struct{}

// MetadataResponse ...
type MetadataResponse struct {
	Metadata metadata.MSystem
}

// MetadataHeadParams ...
type MetadataHeadParams struct {
	ID string `json:"id"`
}

// Head ....
func (ms MetadataService) Head(r *http.Request, args *MetadataHeadParams, reply *MetadataResponse) error {
	log.Printf("metadat head params: %v\n", args)
	s := sessions.Open(args.ID)
	ctx := context.Background()
	return s.Exec(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		head, err := head(r, ctx, u.GetMetadata)
		if err != nil {
			return err
		}
		reply.Metadata = *head
		return err
	})
}

// MetadataGetParams ...
type MetadataGetParams struct {
	ID         string `json:"id"`
	Extraction string // (|STANDARD-XML|COMPACT|COMPACT-INCREMENTAL) the format to pull from the server
}

// Get ....
func (ms MetadataService) Get(r *http.Request, args *MetadataGetParams, reply *MetadataResponse) error {
	log.Printf("metadata get params: %v\n", args)

	s := sessions.Open(args.ID)
	if s == nil {
		return fmt.Errorf("no source found for %s", args.ID)
	}
	log.Printf("connections params for %s %v\n", args.ID, s.Connection)
	if JSONExist(s.MSystem()) {
		standard := metadata.MSystem{}
		JSONLoad(s.MSystem(), &standard)
		reply.Metadata = standard
		return nil
	}
	// lookup the operation for pulling metadata
	if args.Extraction == "" {
		args.Extraction = "COMPACT"
	}
	op, ok := options[args.Extraction]
	if !ok {
		return fmt.Errorf("%s not supported", args.Extraction)
	}
	ctx := context.Background()
	return s.Exec(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		standard, err := op(r, ctx, u.GetMetadata)
		reply.Metadata = *standard
		// bg this
		go func() {
			JSONStore(s.MSystem(), &standard)
		}()
		return err
	})
}

// MetadataRequestType is a typedef metadata extraction options
type MetadataRequestType func(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error)

// options for extracting metadata
var options = map[string]MetadataRequestType{
	"STANDARD-XML":        fullViaStandard,
	"COMPACT":             fullViaCompact,
	"COMPACT-INCREMENTAL": fullViaCompactIncremental,
}

// TODO extract a common func and switch on the incoming param

// fullViaCompactIncremental retrieve the RETS Compact metadata from the server
func fullViaCompactIncremental(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
	compact := &util.IncrementalCompact{}
	err := compact.Load(requester, ctx, url)
	if err != nil {
		return nil, err
	}
	return util.AsStandard(*compact).Convert()
}

// fullViaCompact retrieve the RETS Compact metadata from the server
func fullViaCompact(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
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
	return util.AsStandard(*compact).Convert()
}

// fullViaStandard ...
func fullViaStandard(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
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

// head ...
func head(requester rets.Requester, ctx context.Context, url string) (*metadata.MSystem, error) {
	reader, err := rets.MetadataStream(requester, ctx, rets.MetadataRequest{
		URL:    url,
		Format: "COMPACT",
		MType:  "METADATA-SYSTEM",
		ID:     "0",
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
