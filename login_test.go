/**
	parsing the 'login' action from RETS
 */
package gorets

import (
	"testing"
)

func TestParseCapabilities(t *testing.T) {
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
	urls, err := parseCapability([]byte(body))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}

	if urls.Response.ReplyText != "V2.7.0 2315: Success" {
		t.Errorf("wrong reply code: %s ", urls.Response.ReplyCode)
	}
	if urls.Response.ReplyCode != 0 {
		t.Errorf("wrong reply code: %s ", urls.Response.ReplyCode)
	}
	if urls.Login != "http://server.com:6103/platinum/login" {
		t.Errorf("login urls mismatch: %s ", urls.Login)
	}
	if urls.GetMetadata != "http://server.com:6103/platinum/getmetadata" {
		t.Errorf("login urls mismatch: %s ", urls.Login)
	}
}
