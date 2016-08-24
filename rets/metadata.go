package rets

import (
	"io"
	"net/http"
	"net/url"

	"context"
)

// MetadataRequest ...
type MetadataRequest struct {
	// RETS request options
	URL, HTTPMethod, Format, MType, ID string
}

// MetadataStream ...
func MetadataStream(requester Requester, ctx context.Context, r MetadataRequest) (io.ReadCloser, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	values := url.Query()
	// required
	values.Add("Format", r.Format)
	values.Add("Type", r.MType)
	values.Add("ID", r.ID)

	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}

	url.RawQuery = values.Encode()

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	return DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType)), nil
}
