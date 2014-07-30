/**
provides the photo extraction core testing
*/
package gorets_client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"testing"
)

func TestGetObject(t *testing.T) {
	header := http.Header{}
	textproto.MIMEHeader(header).Add("Content-Type", "image/jpeg")
	textproto.MIMEHeader(header).Add("Content-ID", "123456")
	textproto.MIMEHeader(header).Add("Object-ID", "1")
	textproto.MIMEHeader(header).Add("Preferred", "1")
	textproto.MIMEHeader(header).Add("UID", "1a234234234")
	textproto.MIMEHeader(header).Add("Description", "Outhouse")
	textproto.MIMEHeader(header).Add("Sub-Description", "The urinal")
	textproto.MIMEHeader(header).Add("Location", "http://www.simpleboundary.com/image-5.jpg")

	var body string = `<binary data 1>`
	reader := ioutil.NopCloser(bytes.NewReader([]byte(body)))

	quit := make(chan struct{})
	defer close(quit)
	results := parseGetObjectResult(quit, header, reader)
	result := <-results

	o := result.Object
	ok(t, result.Err)
	equals(t, true, o.Preferred)
	equals(t, "image/jpeg", o.ContentType)
	equals(t, "123456", o.ContentId)
	equals(t, 1, o.ObjectId)
	equals(t, "1a234234234", o.Uid)
	equals(t, "Outhouse", o.Description)
	equals(t, "The urinal", o.SubDescription)
	equals(t, "<binary data 1>", string(o.Blob))
	equals(t, "http://www.simpleboundary.com/image-5.jpg", o.Location)
	equals(t, false, o.RetsError)
}

var boundary string = "simple boundary"

var contentType string = `multipart/parallel; boundary="simple boundary"`

var multipartBody string = `--simple boundary
Content-Type: image/jpeg
Content-ID: 123456
Object-ID: 1
Preferred: 1
ObjectData: ListingKey=123456
ObjectData: ListDate=2013-05-01T12:34:34.8-0500

<binary data 1>
--simple boundary
Content-Type: image/jpeg
Content-ID: 123456
Object-ID: 2
UID: 1a234234234

<binary data 2>
--simple boundary
Content-Type: image/jpeg
Content-ID: 123456
Object-ID: 3
Description: Outhouse
Sub-Description: The urinal

<binary data 3>
--simple boundary
Content-Type: text/xml
Content-ID: 123457
Object-ID: 4
RETS-Error: 1

<RETS ReplyCode="20403" ReplyText="There is no object with that Object-ID"/>

--simple boundary
Content-Type: image/jpeg
Content-ID: 123456
Object-ID: 5
Location: http://www.simpleboundary.com/image-5.jpg

<binary data 5>
--simple boundary--`

func TestExtractBoundary(t *testing.T) {
	extracted := extractBoundary(contentType)

	equals(t, boundary, extracted)
}

func TestGetObjects(t *testing.T) {
	extracted := extractBoundary(contentType)

	equals(t, boundary, extracted)

	body := ioutil.NopCloser(bytes.NewReader([]byte(multipartBody)))

	quit := make(chan struct{})
	defer close(quit)
	results := parseGetObjectsResult(quit, extracted, body)

	r1 := <-results
	ok(t, r1.Err)
	o1 := r1.Object
	equals(t, true, o1.Preferred)
	equals(t, "image/jpeg", o1.ContentType)
	equals(t, "123456", o1.ContentId)
	equals(t, 1, o1.ObjectId)
	equals(t, "<binary data 1>", string(o1.Blob))
	equals(t, "123456", o1.ObjectData["ListingKey"])
	equals(t, "2013-05-01T12:34:34.8-0500", o1.ObjectData["ListDate"])

	r2 := <-results
	ok(t, r2.Err)
	o2 := r2.Object
	equals(t, 2, o2.ObjectId)
	equals(t, "1a234234234", o2.Uid)

	r3 := <-results
	ok(t, r3.Err)
	o3 := r3.Object
	equals(t, 3, o3.ObjectId)
	equals(t, "Outhouse", o3.Description)
	equals(t, "The urinal", o3.SubDescription)

	r4 := <-results
	ok(t, r4.Err)
	o4 := r4.Object
	equals(t, true, o4.RetsError)

	equals(t, "text/xml", o4.ContentType)
	equals(t, "There is no object with that Object-ID", o4.RetsErrorMessage.ReplyText)
	equals(t, 20403, o4.RetsErrorMessage.ReplyCode)

	r5 := <-results
	ok(t, r5.Err)
	o5 := r5.Object
	equals(t, "http://www.simpleboundary.com/image-5.jpg", o5.Location)
	equals(t, "image/jpeg", o5.ContentType)
	equals(t, "123456", o5.ContentId)
	equals(t, 5, o5.ObjectId)
	equals(t, "<binary data 5>", string(o5.Blob))

}

func TestParseGetObjectQuit(t *testing.T) {
	extracted := extractBoundary(contentType)

	equals(t, boundary, extracted)

	body := ioutil.NopCloser(bytes.NewReader([]byte(multipartBody)))

	quit := make(chan struct{})
	defer close(quit)
	results := parseGetObjectsResult(quit, extracted, body)

	r1 := <-results
	ok(t, r1.Err)
	assert(t, r1 != GetObjectResult{},"should not be the zerod object")
	o1 := r1.Object
	equals(t, "image/jpeg", o1.ContentType)
	equals(t, "123456", o1.ContentId)
	equals(t, 1, o1.ObjectId)

	quit <- struct{}{}

	// the closed channel will emit a zero'd value of the proper type
	r2 := <-results
	equals(t, r2, GetObjectResult{})
}
