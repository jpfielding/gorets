package client

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
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
	method := s.HTTPMethodDefault
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Execute(ctx, req)
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
	// clear any state
	err = s.Reset()
	if err != nil {
		return nil, err
	}
	return logoutResponse, nil
}

func processResponseBody(body string) (*LogoutResponse, error) {
	type xmlRets struct {
		XMLName   xml.Name `xml:"RETS"`
		ReplyCode int      `xml:"ReplyCode,attr"`
		ReplyText string   `xml:"ReplyText,attr"`
		Response  string   `xml:"RETS-RESPONSE"`
	}

	rets := xmlRets{}
	decoder := GetXMLReader(bytes.NewBufferString(body), false)
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
		return (&LogoutResponse{
			ReplyCode:      rets.ReplyCode,
			ReplyText:      rets.ReplyText,
			ConnectTime:    connectTime,
			Billing:        values["billing"],
			SignOffMessage: values["signoffmessage"],
		}), nil
	}

	return (&LogoutResponse{
		ReplyCode:      rets.ReplyCode,
		ReplyText:      rets.ReplyText,
		Billing:        values["billing"],
		SignOffMessage: values["signoffmessage"],
	}), nil
}
