/**
 * Created with IntelliJ IDEA.
 * User: jp
 * Date: 9/20/13
 * Time: 11:57 AM
 * To change this template use File | Settings | File Templates.
 */
package gorets

import (
	"fmt"
	"testing"
)

/* http://golang.org/pkg/testing/ */
func TestXxx(t *testing.T) {
	t.Log("bitches love go")
}

func BenchmarkXxx(b *testing.B) {
	b.Log("bitches are fast")
}

func ExampleHello() {
	fmt.Println("hello example")
}


func TestParseCapabilities(t *testing.T) {
	body :=
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
	<RETS-RESPONSE>
	MemberName=Threewide Corporation
	User=2006973, Association Member Primary:Login:Media Restrictions:Office:RETS:RETS Advanced:RETS Customer:System-MRIS:MDS Access Common:MDS Application Login, 90, TWD1
	Broker=TWD,1
	MetadataVersion=1.12.30
	MinMetadataVersion=1.1.1
	OfficeList=TWD;1
	TimeoutSeconds=1800
	Login=http://cornerstone.mris.com:6103/platinum/login
	Action=http://cornerstone.mris.com:6103/platinum/get?Command=Message
	Search=http://cornerstone.mris.com:6103/platinum/search
	Get=http://cornerstone.mris.com:6103/platinum/get
	GetObject=http://cornerstone.mris.com:6103/platinum/getobject
	Logout=http://cornerstone.mris.com:6103/platinum/logout
	GetMetadata=http://cornerstone.mris.com:6103/platinum/getmetadata
	ChangePassword=http://cornerstone.mris.com:6103/platinum/changepassword
	</RETS-RESPONSE>
	</RETS>`
	urls, err := parse([]byte(body))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}

	if urls.Response.ReplyText != "V2.7.0 2315: Success" {
		t.Errorf("wrong reply code: %s ", urls.Response.ReplyCode)
	}
	if urls.Response.ReplyCode != 0 {
		t.Errorf("wrong reply code: %s ", urls.Response.ReplyCode)
	}
	if urls.Login != "http://cornerstone.mris.com:6103/platinum/login" {
		t.Errorf("login urls mismatch: %s ", urls.Login)
	}
	if urls.GetMetadata != "http://cornerstone.mris.com:6103/platinum/getmetadata" {
		t.Errorf("login urls mismatch: %s ", urls.Login)
	}
}
