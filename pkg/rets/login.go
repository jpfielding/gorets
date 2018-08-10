package rets

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"context"
)

// LoginRequest ...
type LoginRequest struct {
	URL, HTTPMethod string
}

// CapabilityURLs ...
type CapabilityURLs struct {
	Response Response

	MemberName,
	User,
	Broker,
	MetadataVersion,
	MinMetadataVersion string
	OfficeList     []string
	TimeoutSeconds int64
	// X-urls
	AdditionalURLs map[string]string
	// urls for web calls
	Action,
	ChangePassword,
	GetMetadata,
	GetObject,
	Login,
	LoginComplete,
	Logout,
	Search,
	Update,
	PostObject,
	GetPayloadList string
}

// Login ...
func Login(ctx context.Context, requester Requester, r LoginRequest) (*CapabilityURLs, error) {
	method := DefaultHTTPMethod
	if r.HTTPMethod != "" {
		method = r.HTTPMethod
	}
	// create a request that we can apply custom headers to
	req, err := http.NewRequest(method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester(ctx, req)
	if err != nil {
		return nil, err
	}
	body := DefaultReEncodeReader(resp.Body, resp.Header.Get(ContentType))

	urls, err := parseCapability(body, r.URL)
	if err != nil {
		return nil, errors.New("unable to parse capabilities response: " + err.Error())
	}
	return urls, nil
}

func parseCapability(body io.ReadCloser, url string) (*CapabilityURLs, error) {
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
	if strings.TrimSpace(rets.Response) == "" {

		detail := ""
		if rets.ReplyCode != StatusOK || rets.ReplyText != "" {
			detail = fmt.Sprintf(" - %d: %s", rets.ReplyCode, rets.ReplyText)
		}
		return nil, errors.New("failed to read urls" + detail)
	}

	reader := bufio.NewReader(strings.NewReader(rets.Response))
	scanner := bufio.NewScanner(reader)

	c := CapabilityURLs{
		AdditionalURLs: map[string]string{},
	}
	values := map[string]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		kv := strings.Split(line, "=")
		// AdditionalURLs
		if strings.HasPrefix(kv[0], "X-") {
			c.AdditionalURLs[kv[0]] = prependHost(url, kv[1])
		}
		// force it to lower case so we can find them in the map
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		value := strings.TrimSpace(kv[1])
		values[key] = value
	}

	c.Action = prependHost(url, values["action"])
	c.ChangePassword = prependHost(url, values["changepassword"])
	c.GetMetadata = prependHost(url, values["getmetadata"])
	c.GetObject = prependHost(url, values["getobject"])
	c.Login = prependHost(url, values["login"])
	c.LoginComplete = prependHost(url, values["logincomplete"])
	c.Logout = prependHost(url, values["logout"])
	c.Search = prependHost(url, values["search"])
	c.Update = prependHost(url, values["update"])
	c.PostObject = prependHost(url, values["postobject"])
	c.GetPayloadList = prependHost(url, values["getpayloadlist"])

	c.TimeoutSeconds, _ = strconv.ParseInt(values["timeoutseconds"], 10, strconv.IntSize)
	c.Response.Code = rets.ReplyCode
	c.Response.Text = rets.ReplyText

	c.MemberName = values["membername"]
	c.User = values["user"]
	c.Broker = values["broker"]
	c.MetadataVersion = values["metadataversion"]
	c.MinMetadataVersion = values["minmetadataversion"]
	c.OfficeList = strings.Split(values["officelist"], ",")

	return &c, nil
}

func prependHost(login, other string) string {
	if other == "" {
		return ""
	}
	otherURL, err := url.Parse(other)
	// todo do something with this err or kill it
	if err != nil {
		return other
	}
	if otherURL.Host != "" {
		return other
	}

	// TODO log or throw
	loginURL, _ := url.Parse(login)
	loginURL.Path = otherURL.Path

	return loginURL.String()
}
