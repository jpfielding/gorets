/**
	provides the photo extraction core
 */
package gorets

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"net/textproto"
	"mime/multipart"
	"strconv"
	"strings"
)

/* 5.5 spec */
type GetObject struct {
	/** required-ish */
	Uid,
	ContentId,
	ContentType string
	/* 5.5.2 this is probably a bad idea, though its solid with the spec */
	ObjectId int
	/** optional */
	Description,
	SubDescription,
	Location string
	/* 5.6.7 - because why would you want to use standard http errors when we can reinvent! */
	RetsError bool
	RetsErrorMessage RetsResponse
	/* 5.6.3 */
	Preferred bool
	/* 5.6.5 */
	ObjectData map[string]string

	Blob []byte
}

type GetObjectResult struct {
	Objects <-chan GetObject
	ProcessingFailure error
}

type GetObjectRequest struct {
	/* 5.3 */
	Url,
	Resource,
	Type,
	Uid,
	Id,
	ObjectId,
	/* TODO support extracting these fields back from the server 5.4.2 */
	ObjectData string
	/* 5.4.1 */
	Location int
}

/*
	GET /platinum/search?
	Class=ALL&
	Count=1&
	Format=COMPACT-DECODED&
	Limit=10&
	Offset=50&
	Query=%28%28LocaleListingStatus%3D%7CACTIVE-CORE%2CCNTG%2FKO-CORE%2CCNTG%2FNO+KO-CORE%2CAPP+REG-CORE%29%2C%7E%28VOWList%3D0%29%29&
	QueryType=DMQL2&
	SearchType=Property
 */
func (s *Session) GetObject(r GetObjectRequest) (*GetObjectResult, error) {
	// required
	values := url.Values{}
	values.Add("Resource", r.Resource)
	values.Add("Type", r.Type)
	values.Add("UID", r.Uid)
	values.Add("object-id", r.ObjectId)

	// optional
	optionalString := func (name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}

	// one or the other _MUST_ be present
	optionalString("ID", r.Id)
	optionalString("UID", r.Uid)

	optionalInt := func (name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d",value))
		}
	}
	optionalInt("Location", r.Location)

	req, err := http.NewRequest(s.HttpMethod, fmt.Sprintf("%s?%s",r.Url,values.Encode()), nil)
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
		return parseGetObjectResult(resp.Header, resp.Body)
	}

	return parseGetObjectsResult(boundary, resp.Body)
}

func parseGetObjectResult(header http.Header, body io.ReadCloser) (*GetObjectResult,error) {
	data := make(chan GetObject)
	result := GetObjectResult{
		Objects: data,
	}
	go func() {
		defer body.Close()
		defer close(data)
		object, err := parseHeadersAndStream(textproto.MIMEHeader(header), body)
		if err != nil {
			result.ProcessingFailure = err
			return
		}
		data <- *object
	}()
	return &result, nil
}

func parseGetObjectsResult(boundary string, body io.ReadCloser) (*GetObjectResult,error) {
	data := make(chan GetObject)
	result := GetObjectResult{
		 Objects: data,
	}
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
				result.ProcessingFailure = err
				return
			}
			object, err := parseHeadersAndStream(part.Header, part)
			if err != nil {
				result.ProcessingFailure = err
				return
			}

			data <- *object
		}
	} ()

	return &result, nil
}

/** TODO - this is the lazy mans version, this needs to be addressed properly */
func extractBoundary(header string) (string) {
	for _,part := range strings.Split(header,";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "boundary=") {
			val := strings.SplitAfterN(part, "=", 2)[1]
			return strings.Trim(val,"\"")
		}
	}
	return ""
}

func parseHeadersAndStream(header textproto.MIMEHeader, body io.ReadCloser) (*GetObject, error) {
	objectId, err := strconv.ParseInt(header.Get("object-id"),10, 64)
	if err != nil {
		return nil, err
	}
	retsError, err := strconv.ParseBool(header.Get("RETS-Error"))
	retsErrorMsg := RetsResponse{}
	switch {
	case err != nil:
		retsError = false
	case retsError:
		// TODO deal with rets error content
	}
	preferred, err := strconv.ParseBool(header.Get("Preferred"))
	if err != nil {
		preferred = false
	}
	objectData := make(map[string]string)
	// TODO extract object data header values
	blob, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	object := GetObject{
		Uid: header.Get("UID"),
		ContentId: header.Get("Content-ID"),
		ContentType: header.Get("Content-Type"),
		ObjectId: int(objectId),
		Description: header.Get("Description"),
		SubDescription: header.Get("Sub-Description"),
		Location: header.Get("Location"),
		RetsError: retsError,
		RetsErrorMessage: retsErrorMsg,
		Preferred: preferred,
		ObjectData: objectData,
		Blob: blob,
	}

	return &object, nil
}
