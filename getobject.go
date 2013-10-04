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
		return parseGetObjectResult(resp.Body)
	}

	return parseGetObjectsResult(boundary, resp.Body)
}

/** TODO - this is the lazy mans version, this needs to be addressed properly */
func extractBoundary(header string) (string) {
	for _,part := range strings.Split(header,";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "boundary=") {
			val := strings.Split(part, "=")[1]
			return strings.Trim(val,"\"")
		}
	}
	return ""
}

func parseGetObjectResult(body io.ReadCloser) (*GetObjectResult,error) {
//	data := make(chan []string)
//	_ <- data

	return nil, nil
}

func parseGetObjectsResult(boundary string, body io.ReadCloser) (*GetObjectResult,error) {
	data := make(chan GetObject)
	result := GetObjectResult{
		 Objects: data,
	}
	go func() {
		partsReader := multipart.NewReader(body, boundary)
		defer body.Close()
		defer close(data)
		for {
			part, err := partsReader.NextPart()
			if err != nil {
				result.ProcessingFailure = err
				return
			}
			objectId, err := strconv.ParseInt(part.Header.Get("object-id"),10, 64)
			if err != nil {
				result.ProcessingFailure = err
				return
			}
			retsError, err := strconv.ParseBool(part.Header.Get("RETS-Error"))
			retsErrorMsg := RetsResponse{}
			switch {
			case err != nil:
					retsError = false
			case retsError:
				// TODO deal with rets error content
			}
			preferred, err := strconv.ParseBool(part.Header.Get("Preferred"))
			if err != nil {
				preferred = false
			}
			objectData := make(map[string]string)
			// TODO extract object data header values
			blob, err := ioutil.ReadAll(part)
			if err != nil {
				result.ProcessingFailure = err
				return
			}

			object := GetObject{
				Uid: part.Header.Get("UID"),
				ContentId: part.Header.Get("Content-ID"),
				ContentType: part.Header.Get("Content-Type"),
				ObjectId: int(objectId),
				Description: part.Header.Get("Description"),
				SubDescription: part.Header.Get("Sub-Description"),
				Location: part.Header.Get("Location"),
				RetsError: retsError,
				RetsErrorMessage: retsErrorMsg,
				Preferred: preferred,
				ObjectData: objectData,
				Blob: blob,
			}
			data <- object
		}
	} ()

	return &result, nil
}
