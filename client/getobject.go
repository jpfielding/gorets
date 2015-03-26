/**
provides the photo extraction core
*/
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
)

/* 5.5 spec */
type GetObject struct {
	/** required */
	ContentId,
	ContentType string
	/* 5.5.2 this is probably a bad idea, though its solid with the spec */
	ObjectId int
	/** optional-ish _must_ return if the request used this field */
	Uid string
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

// helper to abstract the location concept (not thread safe)
func (obj *GetObject) Content() ([]byte, error) {
	if obj == nil {
		return nil, nil
	}
	if obj.Blob != nil {
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

type GetObjectResult struct {
	Object *GetObject
	Err    error
}

type GetObjectRequest struct {
	/* 5.3 */
	URL, HTTPMethod,
	Resource,
	Type,
	Uid,
	/** listing1:1:3:5,listing2:*,listing3:0 */
	Id string
	/** 5.4.2 listing data to be embedded in the response */
	ObjectData []string
	/* 5.4.1 */
	Location int
}

/* */
func (s *Session) GetObject(quit <-chan struct{}, r GetObjectRequest) (<-chan GetObjectResult, error) {
	// required
	values := url.Values{}
	values.Add("Resource", r.Resource)
	values.Add("Type", r.Type)

	// optional
	optionalString := OptionalStringValue(values)

	// one or the other _MUST_ be present
	optionalString("ID", r.Id)
	optionalString("UID", r.Uid)
	// truly optional
	optionalString("ObjectData", strings.Join(r.ObjectData, ","))

	optionalInt := OptionalIntValue(values)
	optionalInt("Location", r.Location)

	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	// TODO use a URL object then properly append to it
	req, err := http.NewRequest(method, fmt.Sprintf("%s?%s", r.URL, values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-type")
	boundary := extractBoundary(contentType)
	if boundary == "" {
		return parseGetObjectResult(quit, resp.Header, resp.Body), nil
	}

	return parseGetObjectsResult(quit, boundary, resp.Body), nil
}

func parseGetObjectResult(quit <-chan struct{}, header http.Header, body io.ReadCloser) <-chan GetObjectResult {
	data := make(chan GetObjectResult)
	go func() {
		defer body.Close()
		defer close(data)
		select {
		case data <- parseHeadersAndStream(textproto.MIMEHeader(header), body):
		case <-quit:
			return
		}
	}()
	return data
}

func parseGetObjectsResult(quit <-chan struct{}, boundary string, body io.ReadCloser) <-chan GetObjectResult {
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
			case <-quit:
				return
			}
		}
	}()

	return data
}

/** TODO - this is the lazy mans version, this needs to be addressed properly */
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
	objectId, err := strconv.ParseInt(header.Get("Object-ID"), 10, 64)
	if err != nil {
		return GetObjectResult{nil, err}
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
		ObjectId:    int(objectId),
		ContentId:   header.Get("Content-ID"),
		ContentType: header.Get("Content-Type"),
		// optional
		Uid:              header.Get("UID"),
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
