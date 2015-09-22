package client

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// LogoutRequest ...
type LogoutRequest struct {
	URL, HTTPMethod string
}

// LogoutResponse ...
type LogoutResponse struct {
	ReplyCode      int
	ReplyText      string
	ConnectTime    uint64
	Billing        string
	SignOffMessage string
}

// Logout ...
func (s *Session) Logout(ctx context.Context, r LogoutRequest) (*LogoutResponse, error) {
	method := "GET"
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, &s.Client, req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	logoutResponse, err := processResponseBody(string(body))
	if err != nil {
		return nil, err
	}

	// wipe the cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	s.Client.Jar = jar
	return logoutResponse, nil
}

func processResponseBody(body string) (*LogoutResponse, error) {
	type XmlRets struct {
		XMLName   xml.Name `xml:"RETS"`
		ReplyCode int      `xml:"ReplyCode,attr"`
		ReplyText string   `xml:"ReplyText,attr"`
		Response  string   `xml:"RETS-RESPONSE"`
	}

	rets := XmlRets{}
	decoder := GetXmlReader(bytes.NewBufferString(body), false)
	err := decoder.Decode(&rets)
	if err != nil && err != io.EOF {
		return nil, err
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

	if val, ok := values["connecttime"]; ok {
		connectTime, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return (&LogoutResponse{rets.ReplyCode, rets.ReplyText, connectTime, values["billing"], values["signoffmessage"]}), nil
	}

	return (&LogoutResponse{ReplyCode: rets.ReplyCode, ReplyText: rets.ReplyText, Billing: values["billing"], SignOffMessage: values["signoffmessage"]}), nil
}
