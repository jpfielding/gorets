package retsutil

import (
	"context"
	"errors"
	"strings"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
)

// IncrementalCompact loads rets.CompcatMetadata in pieces
type IncrementalCompact rets.CompactMetadata

type cmgetter func(string, string) (*rets.CompactMetadata, error)

// GetCompactIncrementalMetadata retrieve the RETS Compact metadata from the server
func (ic *IncrementalCompact) Load(sess rets.Requester, ctx context.Context, url string) error {
	// extract an id'd subesection of metadata
	get := func(id, mtype string) (*rets.CompactMetadata, error) {
		if id == "" {
			id = "0"
		}
		req := rets.MetadataRequest{
			URL:    url,
			Format: "COMPACT",
			MType:  mtype,
			ID:     id,
		}
		reader, er := rets.MetadataStream(sess, ctx, req)
		if er != nil {
			return nil, er
		}
		return rets.ParseMetadataCompactResult(reader)
	}
	msys, err := get("", metadata.MISystem.Name)
	if err != nil {
		return err
	}
	ic.Response = msys.Response
	ic.MSystem = msys.MSystem
	ic.Elements = map[string][]rets.CompactData{}
	// fmt.Printf("compact system: %v\n", compact)
	cds, err := ic.extractChildren(get, []string{}, metadata.MISystem)
	if err != nil {
		return err
	}
	for _, cd := range cds {
		// extract children from this element and put in the system elem
		ic.Elements[cd.Element] = append(ic.Elements[cd.Element], cd)
	}
	return nil
}

func (ic *IncrementalCompact) extractChildren(get cmgetter, idPath []string, mi metadata.MetaInfo) ([]rets.CompactData, error) {
	var tmp []rets.CompactData
	for _, child := range mi.Child {
		cm, err := get(strings.Join(idPath, ":"), child.Name)
		if err != nil {
			return tmp, err
		}
		if cm.Response.Code != rets.StatusOK {
			return tmp, errors.New(cm.Response.Text)
		}
		for _, cdata := range cm.Elements[child.Name] {
			tmp = append(tmp, cdata)
			//  recurse on each member of this cdata
			for _, each := range cdata.Entries() {
				var data map[string]string = each
				idPath := append(idPath, child.ID(data))
				// fmt.Printf("compact system: %v\n", compact)
				cds, err := ic.extractChildren(get, idPath, child)
				if err != nil {
					return tmp, err
				}
				for _, cd := range cds {
					// extract children from this element and put in the system elem
					tmp = append(tmp, cd)
				}
			}
		}
	}
	return tmp, nil
}
