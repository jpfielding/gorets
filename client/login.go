package client

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

// LoginRequest ...
type LoginRequest struct {
	URL, HTTPMethod string
}

// CapabilityURLs ...
type CapabilityURLs struct {
	Response RetsResponse

	MemberName, User, Broker, MetadataVersion, MinMetadataVersion string
	OfficeList                                                    []string
	TimeoutSeconds                                                int64
	// urls for web calls
	Login, Action, Search, Get, GetObject, Logout, GetMetadata, ChangePassword string
}

// Login ...
func Login(requester Requester, ctx context.Context, r LoginRequest) (*CapabilityURLs, error) {
	method := "GET"
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
		return nil, errors.New("unable to parse capabilites response: " + err.Error())
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

	c := CapabilityURLs{}
	c.Login = prependHost(url, values["login"])
	c.Action = prependHost(url, values["action"])
	c.Search = prependHost(url, values["search"])
	c.Get = prependHost(url, values["get"])
	c.GetObject = prependHost(url, values["getobject"])
	c.Logout = prependHost(url, values["logout"])
	c.GetMetadata = prependHost(url, values["getmetadata"])
	c.ChangePassword = prependHost(url, values["changepassword"])

	c.TimeoutSeconds, _ = strconv.ParseInt(values["timeoutseconds"], 10, strconv.IntSize)
	c.Response.ReplyCode = rets.ReplyCode
	c.Response.ReplyText = rets.ReplyText

	c.MemberName = values["membername"]
	c.User = values["user"]
	c.Broker = values["broker"]
	c.MetadataVersion = values["metadataversion"]
	c.MinMetadataVersion = values["minmetadataversion"]
	c.OfficeList = strings.Split(values["officelist"], ",")

	return &c, nil
}

func prependHost(login, other string) string {
	otherURL, err := url.Parse(other)
	// todo do something with this err or kill it
	if err != nil {
		return other
	}
	if otherURL.Host != "" {
		return other
	}

	loginURL, err := url.Parse(login)
	loginURL.Path = otherURL.Path

	return loginURL.String()
}
