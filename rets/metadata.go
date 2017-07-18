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

// PrepMetadataRequest creates an http.Request from a MetadataRequest
func PrepMetadataRequest(r MetadataRequest) (*http.Request, error) {
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

	return http.NewRequest(method, url.String(), nil)
}

// MetadataStream ...
func MetadataStream(ctx context.Context, requester Requester, r MetadataRequest) (io.ReadCloser, error) {
	req, err := PrepMetadataRequest(r)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	return DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType)), nil
}
