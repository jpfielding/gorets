package serverPANIC

import (
	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
)

type Session struct {
	WireLogFile string
	Metadata    metadata.MSystem // probably dont need to cache this in the long term
	// Requester is user state
	Requester rets.Requester
	// URLs need this to know where to route requests
	URLs rets.CapabilityURLs
}
