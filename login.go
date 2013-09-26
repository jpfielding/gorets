/**

 */
package gorets

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
)

type CapabilityUrls struct {
	Response RetsResponse

	MemberName, User, Broker, MetadataVersion, MinMetadataVersion string
	OfficeList []string
	TimeoutSeconds int64
	/** urls for web calls */
	Login,Action,Search,Get,GetObject,Logout,GetMetadata,ChangePassword string
}

func (s *Session) Login(url string) (*CapabilityUrls, error) {
	req, err := s.createRequest(url)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

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
	}
	if s.UserAgentPassword != "" {
		requestId := resp.Header.Get(RETS_REQUEST_ID)
		sessionId := ""// TODO resp.Cookies().Get(RETS_SESSION_ID)
		uaAuthHeader := CalculateUaAuthHeader(s.UserAgent, s.UserAgentPassword, requestId, sessionId, s.Version)
		req.Header.Add(RETS_UA_AUTH_HEADER, uaAuthHeader)
	}

	capabilities, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var urls = &CapabilityUrls{}
	err = urls.parse(capabilities)
	if err != nil {
		return nil, errors.New("unable to parse capabilites response: "+string(capabilities))
	}
	return urls, nil
}



func (c *CapabilityUrls) parse(response []byte) (error){
	type Rets struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		Response string `xml:"RETS-RESPONSE"`
	}

	rets := Rets{}
	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false
	//decoder.AutoClose = append(decoder.AutoClose,"RETS")
	err := decoder.Decode(&rets)
	if err != nil && err != io.EOF {
		return err
	}
	if rets.Response == "" {
		return errors.New("failed to read urls")
	}

	reader := bufio.NewReader(strings.NewReader(rets.Response))
	scanner := bufio.NewScanner(reader)

	values := map[string]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		kv := strings.Split(line, "=")
		// force it to lower case so we can find them in the map
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		value := strings.TrimSpace(kv[1])
		values[key] = value
	}

	c.Login = values["login"]
	c.Action = values["action"]
	c.Search = values["search"]
	c.Get = values["get"]
	c.GetObject = values["getobject"]
	c.Logout = values["logout"]
	c.GetMetadata = values["getmetadata"]
	c.ChangePassword = values["changepassword"]

	c.TimeoutSeconds,_= strconv.ParseInt(values["timeoutseconds"],10,strconv.IntSize)
	c.Response.ReplyCode = rets.ReplyCode
	c.Response.ReplyText = rets.ReplyText

	c.MemberName = values["membername"]
	c.User = values["user"]
	c.Broker = values["broker"]
	c.MetadataVersion = values["metadataversion"]
	c.MinMetadataVersion = values["minmetadataversion"]
	c.OfficeList = strings.Split(values["officelist"],",")

	return nil
}

