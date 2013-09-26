/**
 */
package gorets

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	session.UserAgent = userAgent
	session.HttpMethod = "GET"
	session.Accept = "*/*"
	session.Client = &http.Client{}
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

	Urls *CapabilityUrls
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
	/** urls for web calls */
	Login,Action,Search,Get,GetObject,Logout,GetMetadata,ChangePassword string
}

func (s *Session) Login(loginUrl string) (*CapabilityUrls, error) {
	req, err := http.NewRequest(s.HttpMethod, loginUrl, nil)

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	// TODO do all of this per request
	req.Header.Add(USER_AGENT, s.UserAgent)
	req.Header.Add(RETS_VERSION, s.Version)
	req.Header.Add(ACCEPT, s.Accept)

	fmt.Println("REQUEST:", req)

	resp, err := s.Client.Do(req)

	if err != nil {
		return nil, err
	}

	fmt.Println("RESPONSE:",resp.Header)

	// TODO set digest auth up as a handler
	if resp.StatusCode == 401 {
		challenge := resp.Header.Get(WWW_AUTH)
		if !strings.HasPrefix(strings.ToLower(challenge), "digest") {
			return nil, errors.New("unknown authentication challenge: "+challenge)
		}
		req.Header.Add(WWW_AUTH_RESP, DigestResponse(challenge, s.Username, s.Password, req.Method, req.URL.Path))
		resp, err = s.Client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 401 {
			return nil, errors.New("authentication failed: "+s.Username)
		}
		fmt.Println("RESPONSE (AUTH):",req)
	}
	if s.UserAgentPassword != "" {
		requestId := resp.Header.Get(RETS_REQUEST_ID)
		sessionId := ""// resp.Cookies().Get(RETS_SESSION_ID)
		uaAuthHeader := CalculateUaAuthHeader(s.UserAgent, s.UserAgentPassword, requestId, sessionId, s.Version)
		req.Header.Add(RETS_UA_AUTH_HEADER, uaAuthHeader)
	}

	capabilities, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(capabilities))


	var urls = &CapabilityUrls{}
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
	s.Urls = urls
	return urls, nil
}

