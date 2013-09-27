/**

 */
package gorets

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
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
	req, err := http.NewRequest(s.HttpMethod, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	capabilities, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	urls, err := parse(capabilities)
	if err != nil {
		return nil, errors.New("unable to parse capabilites response: "+string(capabilities))
	}
	return urls, nil
}



func parse(response []byte) (*CapabilityUrls, error){
	type Rets struct {
		XMLName xml.Name `xml:"RETS"`
		ReplyCode int `xml:"ReplyCode,attr"`
		ReplyText string `xml:"ReplyText,attr"`
		Response string `xml:"RETS-RESPONSE"`
	}

	rets := Rets{}
	decoder := xml.NewDecoder(bytes.NewBuffer(response))
	decoder.Strict = false
	err := decoder.Decode(&rets)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if rets.Response == "" {
		return nil, errors.New("failed to read urls")
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

	c := CapabilityUrls{}
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

	return &c, nil
}

