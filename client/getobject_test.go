/**
provides the photo extraction core testing
*/
package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"testing"
	"time"

	testutils "github.com/jpfielding/gorets/testutils"
	"golang.org/x/net/context"
)

func TestGetObject(t *testing.T) {
	header := http.Header{}
	textproto.MIMEHeader(header).Add("Content-Type", "image/jpeg")
	textproto.MIMEHeader(header).Add("Content-ID", "123456")
	textproto.MIMEHeader(header).Add("Object-ID", "1")
	textproto.MIMEHeader(header).Add("Preferred", "1")
	textproto.MIMEHeader(header).Add("UID", "1a234234234")
	textproto.MIMEHeader(header).Add("Content-Description", "Outhouse")
	textproto.MIMEHeader(header).Add("Content-Sub-Description", "The urinal")
	textproto.MIMEHeader(header).Add("Location", "http://www.simpleboundary.com/image-5.jpg")

	var body = `<binary data 1>`
	reader := ioutil.NopCloser(bytes.NewReader([]byte(body)))

	results := parseGetObjectResult(context.Background(), header, reader)
	result := <-results

	o := result.Object
	testutils.Ok(t, result.Err)
	testutils.Equals(t, true, o.Preferred)
	testutils.Equals(t, "image/jpeg", o.ContentType)
	testutils.Equals(t, "123456", o.ContentID)
	testutils.Equals(t, 1, o.ObjectID)
	testutils.Equals(t, "1a234234234", o.UID)
	testutils.Equals(t, "Outhouse", o.Description)
	testutils.Equals(t, "The urinal", o.SubDescription)
	testutils.Equals(t, "<binary data 1>", string(o.Blob))
	testutils.Equals(t, "http://www.simpleboundary.com/image-5.jpg", o.Location)
	testutils.Equals(t, false, o.RetsError)
}

var boundary = "simple boundary"

var contentType = `multipart/parallel; boundary="simple boundary"`

var multipartBody = `--simple boundary
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
Content-Description: Outhouse
Content-Sub-Description: The urinal

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

	testutils.Equals(t, boundary, extracted)
}

func TestGetObjects(t *testing.T) {
	extracted := extractBoundary(contentType)

	testutils.Equals(t, boundary, extracted)

	body := ioutil.NopCloser(bytes.NewReader([]byte(multipartBody)))

	results := parseGetObjectsResult(context.Background(), extracted, body)

	r1 := <-results
	testutils.Ok(t, r1.Err)
	o1 := r1.Object
	testutils.Equals(t, true, o1.Preferred)
	testutils.Equals(t, "image/jpeg", o1.ContentType)
	testutils.Equals(t, "123456", o1.ContentID)
	testutils.Equals(t, 1, o1.ObjectID)
	testutils.Equals(t, "<binary data 1>", string(o1.Blob))
	testutils.Equals(t, "123456", o1.ObjectData["ListingKey"])
	testutils.Equals(t, "2013-05-01T12:34:34.8-0500", o1.ObjectData["ListDate"])

	r2 := <-results
	testutils.Ok(t, r2.Err)
	o2 := r2.Object
	testutils.Equals(t, 2, o2.ObjectID)
	testutils.Equals(t, "1a234234234", o2.UID)

	r3 := <-results
	testutils.Ok(t, r3.Err)
	o3 := r3.Object
	testutils.Equals(t, 3, o3.ObjectID)
	testutils.Equals(t, "Outhouse", o3.Description)
	testutils.Equals(t, "The urinal", o3.SubDescription)

	r4 := <-results
	testutils.Ok(t, r4.Err)
	o4 := r4.Object
	testutils.Equals(t, true, o4.RetsError)

	testutils.Equals(t, "text/xml", o4.ContentType)
	testutils.Equals(t, "There is no object with that Object-ID", o4.RetsErrorMessage.ReplyText)
	testutils.Equals(t, 20403, o4.RetsErrorMessage.ReplyCode)

	r5 := <-results
	testutils.Ok(t, r5.Err)
	o5 := r5.Object
	testutils.Equals(t, "http://www.simpleboundary.com/image-5.jpg", o5.Location)
	testutils.Equals(t, "image/jpeg", o5.ContentType)
	testutils.Equals(t, "123456", o5.ContentID)
	testutils.Equals(t, 5, o5.ObjectID)
	testutils.Equals(t, "<binary data 5>", string(o5.Blob))

}

func TestParseGetObjectQuit(t *testing.T) {
	extracted := extractBoundary(contentType)

	testutils.Equals(t, boundary, extracted)

	body := ioutil.NopCloser(bytes.NewReader([]byte(multipartBody)))

	ctx, cancel := context.WithCancel(context.Background())
	results := parseGetObjectsResult(ctx, extracted, body)

	r1 := <-results
	testutils.Ok(t, r1.Err)
	testutils.Assert(t, r1 != GetObjectResult{}, "should not be the zerod object")
	o1 := r1.Object
	testutils.Equals(t, "image/jpeg", o1.ContentType)
	testutils.Equals(t, "123456", o1.ContentID)
	testutils.Equals(t, 1, o1.ObjectID)

	cancel()
	time.Sleep(100 * time.Millisecond) // I don't like this, but it allows time for the done channel to close.

	// the closed channel will emit a zero'd value of the proper type
	r2 := <-results
	testutils.Equals(t, r2, GetObjectResult{})
}
