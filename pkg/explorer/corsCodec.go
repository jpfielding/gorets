package explorer

import (
	"github.com/gorilla/rpc"

	"net/http"
	"strings"
)

// CodecWithCors creates a custom Codec that adds headers to the WriteResponse
func CodecWithCors(corsDomains []string, baseCodec rpc.Codec) rpc.Codec {
	return corsCodec{corsDomains, baseCodec}
}

type corsCodec struct {
	corsDomains []string
	baseCodec   rpc.Codec
}

// NewRequest ...
func (cc corsCodec) NewRequest(req *http.Request) rpc.CodecRequest {
	return corsCodecRequest{cc.corsDomains, cc.baseCodec.NewRequest(req)}
}

type corsCodecRequest struct {
	corsDomains      []string
	baseCodecRequest rpc.CodecRequest
}

// WriteResponse adds headers onto the ResponseWriter then calls the baseCodecRequest WriteResponse
func (ccr corsCodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, methodErr error) error {
	if len(ccr.corsDomains) > 0 {
		w.Header().Add("Access-Control-Allow-Origin", strings.Join(ccr.corsDomains, " "))
	}
	return ccr.baseCodecRequest.WriteResponse(w, reply, methodErr)
}

func (ccr corsCodecRequest) Method() (string, error) {
	return ccr.baseCodecRequest.Method()
}

func (ccr corsCodecRequest) ReadRequest(args interface{}) error {
	return ccr.baseCodecRequest.ReadRequest(args)
}
