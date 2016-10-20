package explorer

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gorets/retsutil"
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
	fmt.Printf("metadat head params: %v\n", args)
	c := (&ConnectionService{}).Load()[args.ID]
	ctx := context.Background()
	rqr, urls, err := c.Login(ctx)
	if err != nil {
		return err
	}
	head, err := head(rqr, ctx, urls.GetMetadata)
	if err != nil {
		return err
	}
	reply.Metadata = *head
	return nil
}

// MetadataGetParams ...
type MetadataGetParams struct {
	ID         string `json:"id"`
	Extraction string // (|STANDARD-XML|COMPACT|COMPACT-INCREMENTAL) the format to pull from the server
}

// Get ....
func (ms MetadataService) Get(r *http.Request, args *MetadataGetParams, reply *MetadataResponse) error {
	fmt.Printf("metadata get params: %v\n", args)

	c := (&ConnectionService{}).Load()[args.ID]
	if JSONExist(c.MSystem()) {
		standard := metadata.MSystem{}
		JSONLoad(c.MSystem(), &standard)
		reply.Metadata = standard
		return nil
	}
	op, ok := options[args.Extraction]
	if !ok {
		return fmt.Errorf("%s not supported", args.Extraction)
	}
	if args.Extraction == "" {
		args.Extraction = "COMPACT"
	}
	// lookup the operation for pulling metadata
	ctx := context.Background()
	sess, urls, err := c.Login(ctx)
	if err != nil {
		return err
	}
	standard, err := op(sess, ctx, urls.GetMetadata)
	if err != nil {
		return err
	}
	reply.Metadata = *standard
	JSONStore(c.MSystem(), &standard)
	return nil
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
	compact := &retsutil.IncrementalCompact{}
	err := compact.Load(requester, ctx, url)
	if err != nil {
		return nil, err
	}
	return retsutil.AsStandard(*compact).Convert()
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
	return retsutil.AsStandard(*compact).Convert()
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
