package retsutil

import (
	"context"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
)

// IncrementalCompact loads rets.CompcatMetadata in pieces
type IncrementalCompact rets.CompactMetadata

// GetCompactIncrementalMetadata retrieve the RETS Compact metadata from the server
func (ic *IncrementalCompact) Load(sess rets.Requester, ctx context.Context, url string) error {
	// extract an id'd subesection of metadata
	get := func(id, mtype string) (*rets.CompactMetadata, error) {
		req := rets.MetadataRequest{
			URL:    url,
			Format: "COMPACT",
			MType:  metadata.MISystem.Name,
			ID:     id,
		}
		reader, er := rets.MetadataStream(sess, ctx, req)
		if er != nil {
			return nil, er
		}
		return rets.ParseMetadataCompactResult(reader)
	}
	msys, err := get("0", metadata.MISystem.Name)
	if err != nil {
		return err
	}
	ic.Response = msys.Response
	ic.MSystem = msys.MSystem
	// fmt.Printf("compact system: %v\n", compact)
	return nil
}
