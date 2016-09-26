/**
parsing the 'login' action from RETS
*/
package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
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
	</RETS-RESPONSE>
	</RETS>`
	urls, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	testutils.Ok(t, err)

	testutils.Equals(t, urls.Response.Text, "V2.7.0 2315: Success")
	testutils.Equals(t, urls.Response.Code, StatusOK)
	testutils.Equals(t, urls.Login, "http://server.com:6103/platinum/login")
	testutils.Equals(t, urls.GetMetadata, "http://server.com:6103/platinum/getmetadata")
}

func TestPrependHost(t *testing.T) {
	login := "http://server.com:6103/platinum/login"
	abs := "http://server.com:6103/platinum/search"
	rel := "/platinum/search"
	testutils.Equals(t, abs, prependHost(login, abs))
	testutils.Equals(t, abs, prependHost(login, rel))
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
	</RETS-RESPONSE>
	</RETS>`
	urls, err := parseCapability(
		ioutil.NopCloser(strings.NewReader(body)),
		"http://server.com:6103/platinum/login",
	)
	testutils.Ok(t, err)

	testutils.Equals(t, urls.Response.Text, "V2.7.0 2315: Success")
	testutils.Equals(t, urls.Response.Code, StatusOK)
	testutils.Equals(t, urls.Login, "http://server.com:6103/platinum/login")
	testutils.Equals(t, urls.GetMetadata, "http://server.com:6103/platinum/getmetadata")
}
