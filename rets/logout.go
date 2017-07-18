package rets

import (
	"bufio"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"strings"

	"context"
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
func Logout(ctx context.Context, requester Requester, r LogoutRequest) (*LogoutResponse, error) {
	method := DefaultHTTPMethod
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	body := DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType))

	logoutResponse, err := processResponseBody(body)
	if err != nil {
		return nil, err
	}
	// TODO clear any state
	return logoutResponse, nil
}

func processResponseBody(body io.ReadCloser) (*LogoutResponse, error) {
	defer body.Close()
	type xmlRets struct {
		XMLName   xml.Name `xml:"RETS"`
		ReplyCode int      `xml:"ReplyCode,attr"`
		ReplyText string   `xml:"ReplyText,attr"`
		Response  string   `xml:"RETS-RESPONSE"`
	}

	rets := xmlRets{}
	decoder := DefaultXMLDecoder(body, false)
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
