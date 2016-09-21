package rets

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"context"

	"golang.org/x/net/context/ctxhttp"
)

// NewObjectFromStream ...
func NewObjectFromStream(header textproto.MIMEHeader, body io.ReadCloser) (*Object, error) {
	objectID, err := strconv.ParseInt(header.Get("Object-ID"), 10, 64)
	if err != nil {
		// Attempt to parse a Rets Response code (if it exists)
		resp, parseErr := ReadResponse(body)
		if parseErr != nil {
			return nil, err
		}
		// Include a GetObject (empty of content) so that its rets response can be retrieved
		emptyResult := Object{
			RetsMessage: resp,
			RetsError:   resp.Code != 0,
		}
		return &emptyResult, err
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
		return nil, err
	}

	// 5.6.7
	retsError, err := strconv.ParseBool(header.Get("RETS-Error"))
	retsMsg, err := ReadResponse(ioutil.NopCloser(bytes.NewReader(blob)))

	// there is a rets message, stash it and wipe the content
	if err == nil {
		blob = nil
	}

	object := Object{
		// required
		ObjectID:    int(objectID),
		ContentID:   header.Get("Content-ID"),
		ContentType: header.Get("Content-Type"),
		// optional
		UID:            header.Get("UID"),
		Description:    header.Get("Content-Description"),
		SubDescription: header.Get("Content-Sub-Description"),
		Location:       header.Get("Location"),
		RetsError:      retsError,
		RetsMessage:    retsMsg,
		Preferred:      preferred,
		ObjectData:     objectData,
		Blob:           blob,
	}

	return &object, nil
}

// Object provides the photo extraction core for RETS section 5.5
type Object struct {
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
	RetsError   bool
	RetsMessage *Response
	/* 5.6.3 */
	Preferred bool
	/* 5.6.5 */
	ObjectData map[string]string
	/** it may be wiser to convert this to a readcloser with a content-length */
	Blob []byte
}

// ContentWithContext is a helper to abstract the location concept (not thread safe)
func (obj *Object) ContentWithContext(ctx context.Context) ([]byte, error) {
	if obj == nil {
		return nil, nil
	}
	if len(obj.Blob) > 0 || obj.Location == "" {
		return obj.Blob, nil
	}
	resp, err := ctxhttp.Get(ctx, nil, obj.Location)
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

// Content is a helper to abstract the location concept (not thread safe)
func (obj *Object) Content() ([]byte, error) {
	return obj.ContentWithContext(context.Background())
}
