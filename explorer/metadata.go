package explorer

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jpfielding/gorets/config"
	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gorets/util"
)

// MetadataService ...
type MetadataService struct{}

// MetadataResponse ...
type MetadataResponse struct {
	Metadata metadata.MSystem `json:"Metadata"`
	Wirelog  string           `json:"wirelog,omitempty"`
}

// MetadataGetParams ...
type MetadataGetParams struct {
	Connection config.Config `json:"connection"`
	Oldest     time.Duration `json:"oldest"` // the oldst metadata we're willing to accept (minutes)
	Extraction string        // (|STANDARD-XML|COMPACT|COMPACT-INCREMENTAL) the format to pull from the server
}

// Get ....
func (ms MetadataService) Get(r *http.Request, args *MetadataGetParams, reply *MetadataResponse) error {
	fmt.Printf("metadata get params: %v\n", args)

	cfg := args.Connection
	// TOOD make a head request and see if how stale this is??
	if JSONExist(MSystem(cfg), args.Oldest*time.Minute) {
		fmt.Printf("found cached (<%dm old) metadata for %s\n", args.Oldest, cfg.ID)
		standard := metadata.MSystem{}
		JSONLoad(MSystem(cfg), &standard)
		reply.Metadata = standard
		return nil
	}
	// lookup the operation for pulling metadata
	if args.Extraction == "" {
		// TODO deal with sources not supporting the default type
		args.Extraction = "COMPACT"
	}
	op, ok := options[args.Extraction]
	if !ok {
		return fmt.Errorf("%s not supported", args.Extraction)
	}
	ctx := context.Background()
	wirelog := bytes.Buffer{}
	sess, err := cfg.Connect(ctx, &wirelog)
	if err != nil {
		return err
	}
	defer sess.Close()
	return sess.Process(ctx, func(r rets.Requester, u rets.CapabilityURLs) error {
		fmt.Printf("requesting remote metadata for %s\n", cfg.ID)
		standard, err := op(ctx, r, u.GetMetadata)
		if err != nil {
			reply.Wirelog = string(wirelog.Bytes())
			return err
		}
		reply.Metadata = *standard
		reply.Wirelog = string(wirelog.Bytes())
		// bg this
		go func() {
			JSONStore(MSystem(cfg), &standard)
		}()
		return err
	})
}

// MSystem path
func MSystem(c config.Config) string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.ID)
}

// MetadataRequestType is a typedef metadata extraction options
type MetadataRequestType func(ctx context.Context, requester rets.Requester, url string) (*metadata.MSystem, error)

// options for extracting metadata
var options = map[string]MetadataRequestType{
	"STANDARD-XML":        fullViaStandard,
	"COMPACT":             fullViaCompact,
	"COMPACT-INCREMENTAL": fullViaCompactIncremental,
}

// TODO extract a common func and switch on the incoming param

// fullViaCompactIncremental retrieve the RETS Compact metadata from the server
func fullViaCompactIncremental(ctx context.Context, requester rets.Requester, url string) (*metadata.MSystem, error) {
	compact := &util.IncrementalCompact{}
	err := compact.Load(ctx, requester, url)
	if err != nil {
		return nil, err
	}
	return util.AsStandard(*compact).Convert()
}

// fullViaCompact retrieve the RETS Compact metadata from the server
func fullViaCompact(ctx context.Context, requester rets.Requester, url string) (*metadata.MSystem, error) {
	reader, err := rets.MetadataStream(rets.MetadataResponse(ctx, requester, rets.MetadataRequest{
		URL: url,
		MetadataParams: rets.MetadataParams{
			Format: "COMPACT",
			MType:  "METADATA-SYSTEM",
			ID:     "*",
		},
	}))
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
func fullViaStandard(ctx context.Context, requester rets.Requester, url string) (*metadata.MSystem, error) {
	reader, err := rets.MetadataStream(rets.MetadataResponse(ctx, requester, rets.MetadataRequest{
		URL: url,
		MetadataParams: rets.MetadataParams{
			Format: "STANDARD-XML",
			MType:  "METADATA-SYSTEM",
			ID:     "*",
		},
	}))
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
func head(ctx context.Context, requester rets.Requester, url string) (*metadata.MSystem, error) {
	params := rets.MetadataParams{
		Format: "COMPACT",
		MType:  "METADATA-SYSTEM",
		ID:     "0",
	}
	reader, err := rets.MetadataStream(rets.MetadataResponse(ctx, requester, rets.MetadataRequest{
		URL:            url,
		MetadataParams: params,
	}))
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
