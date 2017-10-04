package rets

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"context"
)

// PrepGetObjects creates an http.Request from a GetObjectRequest
func PrepGetObjects(r GetObjectRequest) (*http.Request, error) {
	url, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	values := url.Query()

	// required
	values.Add("Resource", r.Resource)
	values.Add("Type", r.Type)

	// optional
	optionalString := OptionalStringValue(values)

	// one or the other _MUST_ be present
	optionalString("ID", r.ID)
	optionalString("UID", r.UID)
	// truly optional
	optionalString("ObjectData", strings.Join(r.ObjectData, ","))

	optionalInt := OptionalIntValue(values)
	optionalInt("Location", r.Location)

	method := DefaultHTTPMethod
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

// GetObjects sends the GetObject request
func GetObjects(ctx context.Context, requester Requester, r GetObjectRequest) (*http.Response, error) {
	req, err := PrepGetObjects(r)
	if err != nil {
		return nil, err
	}
	return requester(ctx, req)
}

// GetObjectResponse is the response holder for processing GetObject requests
type GetObjectResponse struct {
	Response *http.Response
}

// ForEach ...
func (r *GetObjectResponse) ForEach(result GetObjectResult) error {
	resp := r.Response
	defer resp.Body.Close()
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	// its not multipart, just leave
	if !strings.HasPrefix(mediaType, "multipart/") {
		return result(NewObjectFromStream(textproto.MIMEHeader(resp.Header), resp.Body))
	}
	// its multipart, need to break it up
	partsReader := multipart.NewReader(resp.Body, params["boundary"])
	for {
		part, err := partsReader.NextPart()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}
		err = result(NewObjectFromStream(part.Header, part))
		if err != nil {
			return err
		}
	}
	// return nil
}

// Close ...
func (r *GetObjectResponse) Close() error {
	if r == nil || r.Response.Body == nil {
		return nil
	}
	return r.Response.Body.Close()
}

// GetObjectResult is the callback walking func for retrieving objets
type GetObjectResult func(*Object, error) error

// GetObjectParams holds the parameters for GetObject requests
type GetObjectParams struct {
	Resource,
	Type,
	UID,
	// ID listing1:1:3:5,listing2:*,listing3:0 */
	ID string
	/** 5.4.2 listing data to be embedded in the response */
	ObjectData []string
	/* 5.4.1 */
	Location int
}

// GetObjectRequest ...
type GetObjectRequest struct {
	/* 5.3 */
	URL,
	HTTPMethod string
	HTTPFormEncodedValues bool // POST style http params
	GetObjectParams
}
