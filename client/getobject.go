package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

// GetObject provides the photo extraction core for RETS section 5.5
type GetObject struct {
	// ContentID required
	ContentID,
	ContentType string
	// ObjectID 5.5.2 this is probably a bad idea, though its solid with the spec
	ObjectID int
	// optional-ish _must_ return if the request used this field
	UID string
	/** optional */
	Description,
	SubDescription,
	Location string
	/* 5.6.7 - because why would you want to use standard http errors when we can reinvent! */
	RetsError        bool
	RetsErrorMessage RetsResponse
	/* 5.6.3 */
	Preferred bool
	/* 5.6.5 */
	ObjectData map[string]string
	/** it may be wiser to convert this to a readcloser with a content-length */
	Blob []byte
}

// Content is a helper to abstract the location concept (not thread safe)
func (obj *GetObject) Content() ([]byte, error) {
	if obj == nil {
		return nil, nil
	}
	if len(obj.Blob) > 0 || obj.Location == "" {
		return obj.Blob, nil
	}
	resp, err := http.Get(obj.Location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		obj.ContentType = ct
	}
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	obj.Blob = blob
	return blob, nil
}

// GetObjectResult ...
type GetObjectResult struct {
	Object *GetObject
	Err    error
}

// GetObjectRequest ...
type GetObjectRequest struct {
	/* 5.3 */
	URL, HTTPMethod,
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

// GetObject ...
func (s *Session) GetObject(ctx context.Context, r GetObjectRequest) (<-chan GetObjectResult, error) {
	// required
	values := url.Values{}
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

	method := s.HTTPMethodDefault
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	// TODO use a URL object then properly append to it
	req, err := http.NewRequest(method, fmt.Sprintf("%s?%s", r.URL, values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-type")
	boundary := extractBoundary(contentType)
	if boundary == "" {
		return parseGetObjectResult(ctx, resp.Header, resp.Body), nil
	}

	return parseGetObjectsResult(ctx, boundary, resp.Body), nil
}

func parseGetObjectResult(ctx context.Context, header http.Header, body io.ReadCloser) <-chan GetObjectResult {
	data := make(chan GetObjectResult)
	go func() {
		defer body.Close()
		defer close(data)
		select {
		case data <- parseHeadersAndStream(textproto.MIMEHeader(header), body):
		case <-ctx.Done():
			return
		}
	}()
	return data
}

func parseGetObjectsResult(ctx context.Context, boundary string, body io.ReadCloser) <-chan GetObjectResult {
	data := make(chan GetObjectResult)
	go func() {
		defer body.Close()
		defer close(data)
		partsReader := multipart.NewReader(body, boundary)
		for {
			part, err := partsReader.NextPart()
			switch {
			case err == io.EOF:
				return
			case err != nil:
				data <- GetObjectResult{nil, err}
				return
			}

			select {
			case data <- parseHeadersAndStream(part.Header, part):
			case <-ctx.Done():
				return
			}
		}
	}()

	return data
}

// TODO - this is the lazy mans version, this needs to be addressed properly
func extractBoundary(header string) string {
	for _, part := range strings.Split(header, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "boundary=") {
			val := strings.SplitAfterN(part, "=", 2)[1]
			return strings.Trim(val, "\"")
		}
	}
	return ""
}

func parseHeadersAndStream(header textproto.MIMEHeader, body io.ReadCloser) GetObjectResult {
	objectID, err := strconv.ParseInt(header.Get("Object-ID"), 10, 64)
	if err != nil {
		// Attempt to parse a Rets Response code (if it exists)
		retsResp, parseErr := ParseRetsResponse(body)
		if parseErr != nil {
			return GetObjectResult{nil, err}
		}
		// Include a GetObject (empty of content) so that its rets response can be retrieved
		emptyResult := GetObject{
			RetsErrorMessage: *retsResp,
			RetsError:        retsResp.ReplyCode != 0,
		}
		return GetObjectResult{&emptyResult, err}
	}
	preferred, err := strconv.ParseBool(header.Get("Preferred"))
	if err != nil {
		preferred = false
	}
	objectData := make(map[string]string)
	for _, v := range header[textproto.CanonicalMIMEHeaderKey("ObjectData")] {
		kv := strings.Split(v, "=")
		objectData[kv[0]] = kv[1]
	}
	blob, err := ioutil.ReadAll(body)
	if err != nil {
		return GetObjectResult{nil, err}
	}

	retsError, err := strconv.ParseBool(header.Get("RETS-Error"))
	retsErrorMsg := &RetsResponse{0, "Success"}
	switch {
	case err != nil:
		retsError = false
	case retsError:
		body := ioutil.NopCloser(bytes.NewReader([]byte(blob)))
		retsErrorMsg, err = ParseRetsResponse(body)
		if err != nil {
			return GetObjectResult{nil, err}
		}
	}

	object := GetObject{
		// required
		ObjectID:    int(objectID),
		ContentID:   header.Get("Content-ID"),
		ContentType: header.Get("Content-Type"),
		// optional
		UID:              header.Get("UID"),
		Description:      header.Get("Content-Description"),
		SubDescription:   header.Get("Content-Sub-Description"),
		Location:         header.Get("Location"),
		RetsError:        retsError,
		RetsErrorMessage: *retsErrorMsg,
		Preferred:        preferred,
		ObjectData:       objectData,
		Blob:             blob,
	}

	return GetObjectResult{&object, nil}
}
