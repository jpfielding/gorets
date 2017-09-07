/**
parsing the 'login' action from RETS
*/
package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCapabilitiesAbsoluteUrls(t *testing.T) {
	body :=
		`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
	<RETS-RESPONSE>
	MemberName=Threewide Corporation
	User=2000343, Association Member Primary:Login:Media Restrictions:Office:RETS:RETS Advanced:RETS Customer:System-MRIS:MDS Access Common:MDS Application Login, 90, TURD1
	Broker=TWD,1
	MetadataVersion=1.12.30
	MinMetadataVersion=1.1.1
	OfficeList=TWD;1
	TimeoutSeconds=1800
	Login=http://server.com:6103/platinum/login
	Action=http://server.com:6103/platinum/get?Command=Message
	Search=http://server.com:6103/platinum/search
	Get=http://server.com:6103/platinum/get
	GetObject=http://server.com:6103/platinum/getobject
	Logout=http://server.com:6103/platinum/logout
	GetMetadata=http://server.com:6103/platinum/getmetadata
	ChangePassword=http://server.com:6103/platinum/changepassword
	X-SampleLinks=/rets2_2/Links
	X-SupportSite=http://flexmls.com/rets/
	X-NotificationFeed=http://example.com/atom/feed/private/atom.xml
	</RETS-RESPONSE>
	</RETS>`
	urls, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	assert.Nil(t, err)

	assert.Equal(t, urls.Response.Text, "V2.7.0 2315: Success")
	assert.Equal(t, urls.Response.Code, StatusOK)
	assert.Equal(t, urls.Login, "http://server.com:6103/platinum/login")
	assert.Equal(t, urls.GetMetadata, "http://server.com:6103/platinum/getmetadata")
	assert.Equal(t, "http://example.com/atom/feed/private/atom.xml", urls.AdditionalURLs["X-NotificationFeed"])
}

func TestPrependHost(t *testing.T) {
	login := "http://server.com:6103/platinum/login"
	abs := "http://server.com:6103/platinum/search"
	rel := "/platinum/search"
	assert.Equal(t, abs, prependHost(login, abs))
	assert.Equal(t, abs, prependHost(login, rel))
}

func TestParseCapabilitiesRelativeUrls(t *testing.T) {
	body :=
		`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
	<RETS-RESPONSE>
	MemberName=Threewide Corporation
	User=2000343, Association Member Primary:Login:Media Restrictions:Office:RETS:RETS Advanced:RETS Customer:System-MRIS:MDS Access Common:MDS Application Login, 90, TURD1
	Broker=TWD,1
	MetadataVersion=1.12.30
	MinMetadataVersion=1.1.1
	OfficeList=TWD;1
	TimeoutSeconds=1800
	Login=/platinum/login
	Action=/platinum/get?Command=Message
	Search=/platinum/search
	Get=/platinum/get
	GetObject=/platinum/getobject
	Logout=/platinum/logout
	GetMetadata=/platinum/getmetadata
	ChangePassword=/platinum/changepassword
	X-SampleLinks=/rets2_2/Links
	X-SupportSite=http://flexmls.com/rets/
	X-NotificationFeed=http://example.com/atom/feed/private/atom.xml
	</RETS-RESPONSE>
	</RETS>`
	urls, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	assert.Nil(t, err)

	assert.Equal(t, urls.Response.Text, "V2.7.0 2315: Success")
	assert.Equal(t, urls.Response.Code, StatusOK)
	assert.Equal(t, urls.Login, "http://server.com:6103/platinum/login")
	assert.Equal(t, urls.GetMetadata, "http://server.com:6103/platinum/getmetadata")
	assert.Equal(t, "http://server.com:6103/rets2_2/Links", urls.AdditionalURLs["X-SampleLinks"])
}

func TestFailedLoginNotRets(t *testing.T) {
	body :=
		`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1"/>
	<title>401 - Unauthorized: Access is denied due to invalid credentials.</title>
	</head>
	<body>
	<div id="header"><h1>Server Error</h1></div>
	<div id="content">
	 <div class="content-container"><fieldset>
	  <h2>401 - Unauthorized: Access is denied due to invalid credentials.</h2>
	  <h3>You do not have permission to view this directory or page using the credentials that you supplied.</h3>
	 </fieldset></div>
	</div>
	</body>
	</html>`

	_, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	assert.NotNil(t, err)
	assert.Equal(t, "expected element type <RETS> but have <html>", err.Error())
}

func TestFailedLoginRets(t *testing.T) {
	body :=
		`<RETS>
	<RETS-RESPONSE>
	</RETS-RESPONSE>
	</RETS>`

	_, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to read urls", err.Error())
}

func TestFailedLoginRetsWithDetails(t *testing.T) {
	body :=
		`<RETS ReplyCode="20036" ReplyText="Missing User-Agent request header field." >
	<RETS-RESPONSE>
	</RETS-RESPONSE>
	</RETS>`

	_, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to read urls - 20036: Missing User-Agent request header field.", err.Error())
}
