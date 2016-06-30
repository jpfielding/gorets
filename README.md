gorets
======

RETS client in Go

[![Build Status](https://travis-ci.org/jpfielding/gorets.svg?branch=master)](https://travis-ci.org/jpfielding/gorets.

The attempt is to meet 1.8.0 compliance.

http://www.reso.org/assets/RETS/Specifications/rets_1_8.pdf.

Find me at gophers.slack.com#gorets

```go
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

func main() {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	loginURL := flag.String("login-url", "", "Login URL for the RETS server")
	userAgent := flag.String("user-agent", "Threewide/1.0", "User agent for the RETS client")
	userAgentPw := flag.String("user-agent-pw", "", "User agent password (for RETS UA-Auth)")
	retsVersion := flag.String("rets-version", "", "RETS Version")
	logFile := flag.String("log-file", "", "")

	flag.Parse()

	transport := wirelog.NewHTTPTransport()

	if *logFile != "" {
		wirelog.LogToFile(transport, *logFile, true, true)
	}

	// should we throw an err here too?
	session, err := rets.DefaultSession(
	 	*username,
	 	*password,
	 	*userAgent,
	 	*userAgentPw,
	 	*retsVersion,
	 	&transport,
	 	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	capability, err := rets.Login(session, ctx, rets.LoginRequest{URL: *loginURL})
	if err != nil {
		panic(err)
	}
	fmt.Println("Login: ", capability.Login)
	fmt.Println("Metadata: ", capability.GetMetadata)
	fmt.Println("Search: ", capability.Search)
	fmt.Println("GetObject: ", capability.GetObject)

	err = rets.Get(session, ctx, rets.GetRequest{URL: capability.Get})
	if err != nil {
		fmt.Println("this was stupid, shouldnt even be here")
	}

	mURL := capability.GetMetadata
	format := "COMPACT"
	rets.GetCompactMetadata(session, ctx, rets.MetadataRequest{
		URL:    mURL,
		Format: format,
		MType:  "METADATA-SYSTEM",
		ID:     "0",
	})
	//	session.GetMetadata(rets.MetadataRequest{mUrl, format, "METADATA-RESOURCE", "0"})
	//	session.GetMetadata(rets.MetadataRequest{mUrl, format, "METADATA-CLASS", "ActiveAgent"})
	//	session.GetMetadata(rets.MetadataRequest{mUrl, format, "METADATA-TABLE", "ActiveAgent:ActiveAgent"})

	req := rets.SearchRequest{
		URL:        capability.Search,
		Query:      "((180=|AH))",
		SearchType: "Property",
		Class:      "1",
		Format:     "COMPACT-DECODED",
		QueryType:  "DMQL2",
		Count:      rets.CountIncluded,
		Limit:      3,
		Offset:     -1,
	}
	result, err := rets.SearchCompact(session, ctx, req)
	defer result.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("COLUMNS:", result.Columns)
	hasMoreRows, err := result.ForEach(func(row []string, err error) error {
		fmt.Println(row)
		return nil
	})
	if hasMoreRows {
		fmt.Println("more rows available")
	}
	if err != nil {
		panic(err)
	}

	one, err := rets.GetObject(session, ctx, rets.GetObjectRequest{
		URL:      capability.GetObject,
		Resource: "Property",
		Type:     "Photo",
		Id:       "3986587:1",
	})
	if err != nil {
		panic(err)
	}
	for r := range one {
		if err != nil {
			panic(err)
		}
		o := r.Object
		fmt.Println("PHOTO-META: ", o.ContentType, o.ContentId, o.ObjectId, len(o.Blob))
	}
	all, err := rets.GetObject(session, ctx, rets.GetObjectRequest{
		URL:      capability.GetObject,
		Resource: "Property",
		Type:     "Photo",
		Id:       "3986587:*",
	})
	if err != nil {
		panic(err)
	}
	for r := range all {
		if err != nil {
			panic(err)
		}
		o := r.Object
		fmt.Println("PHOTO-META: ", o.ContentType, o.ContentId, o.ObjectId, len(o.Blob))
	}

	rets.Logout(session, ctx, rets.LogoutRequest{URL: capability.Logout})
}
```
