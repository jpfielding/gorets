/**
 * Created with IntelliJ IDEA.
 * User: jp
 * Date: 9/20/13
 * Time: 11:55 AM
 * To change this template use File | Settings | File Templates.
 */
/*
Package main for the start of a rets client in go
*/
package gorets

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

const DEFAULT_TIMEOUT int = 300000

/** header values */
const (
	RETS_VERSION string = "RETS-Version"
	RETS_SESSION_ID string = "RETS-Session-ID"
	RETS_REQUEST_ID string = "RETS-Request-ID"
	USER_AGENT string = "User-Agent"
	RETS_UA_AUTH_HEADER string = "RETS-UA-Authorization"
	ACCEPT string = "Accept"
	ACCEPT_ENCODING string = "Accept-Encoding"
	CONTENT_ENCODING string = "Content-Encoding"
	DEFLATE_ENCODINGS string = "gzip,deflate"
	CONTENT_TYPE string = "Content-Type"
	WWW_AUTH string = "Www-Authenticate"
	WWW_AUTH_RESP string = "Authorization"

)

func NewSession(user, pw, userAgent, userAgentPw string) *Session {
	var session Session
	session.Username = user
	session.Password = pw
	session.Version = "RETS/1.5"
	session.UserAgent = "Threewide/1.5"
	return &session
}

type Session struct {
	Username,Password string
	UserAgentPassword string
	HttpMethod string
	Version string // "Threewide/1.5"

	Accept string
	UserAgent string

	Client *http.Client

	Urls CapabilityUrls
}


type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

type CapabilityUrls struct {
	Response RetsResponse

	MemberName, User, Broker, MetadataVersion, MinMetadataVersion string
	OfficeList []string
	TimeoutSeconds int
	Login,Action,Search,Get,GetObject,Logout,GetMetadata,ChangePassword string
}

func (r *Session) Login(name string) (*CapabilityUrls, error) {
	req, err := http.NewRequest(r.httpMethod, r.loginUrl, nil)
	if err != nil {
		fmt.Println(err)
		// handle error
	}

	// TODO do all of this per request
	req.Header.Add(USER_AGENT, r.UserAgent)
	req.Header.Add(RETS_VERSION, r.Version)
	req.Header.Add(ACCEPT, r.Accept)

	fmt.Println("REQUEST:", req)

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	fmt.Println("RESPONSE:",resp.Header)

	// TODO set digest auth up as a handler
	if resp.StatusCode == 401 {
		challenge := resp.Header.Get(WWW_AUTH)
		if !strings.HasPrefix(strings.ToLower(challenge), "digest") {
			panic("recognized challenge: "+ challenge)
		}
		req.Header.Add(WWW_AUTH_RESP, auth.DigestResponse(challenge, r.username, r.password, req.Method, req.URL.Path))
		resp, err = r.client.Do(req)
		if err != nil {
			return nil, err
		}
		fmt.Println("RESPONSE (AUTH):",resp)
	}
	if r.UserAgentPassword != nil {
		requestId := resp.Header.Get(RETS_REQUEST_ID)
		sessionId := ""// resp.Cookies().Get(RETS_SESSION_ID)
		uaAuthHeader := auth.CalculateUaAuthHeader(r.UserAgent, r.UserAgentPassword, requestId, sessionId, r.Version)
		req.Header.Add(RETS_UA_AUTH_HEADER, uaAuthHeader)
	}

	capabilities, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(capabilities))


	/*
<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
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
	 */
	return &CapabilityUrls{}
}

