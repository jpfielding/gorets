package rets

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"context"
)

// MetadataParams for the request
type MetadataParams struct {
	Format, MType, ID string
}

// MetadataRequest ...
type MetadataRequest struct {
	// RETS request options
	URL, HTTPMethod       string
	HTTPFormEncodedValues bool // POST style http params
	MetadataParams
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

	// http POST style params
	if r.HTTPFormEncodedValues {
		req, err := http.NewRequest(method, url.String(), bytes.NewBufferString(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, err
	}
	url.RawQuery = values.Encode()
	return http.NewRequest(method, url.String(), nil)
}

// MetadataResponse processes the request... TODO may no longer be necessary
func MetadataResponse(ctx context.Context, requester Requester, r MetadataRequest) (*http.Response, error) {
	req, err := PrepMetadataRequest(r)
	if err != nil {
		return nil, err
	}

	return requester(ctx, req)
}

// MetadataStream encodes the http response stream for us
func MetadataStream(resp *http.Response, err error) (io.ReadCloser, error) {
	return DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType)), nil
}
