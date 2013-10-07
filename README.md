gorets
======

RETS client in Go


The attempt is to meet 1.8.0 compliance.

http://www.reso.org/assets/RETS/Specifications/rets_1_8.pdf.

```go
package main


import (
	"flag"
	"fmt"
	"os"
	"github.com/jpfielding/gorets"
	"io"
)

func main () {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	loginUrl := flag.String("login-url", "", "Login URL for the RETS server")
	userAgent := flag.String("user-agent","Threewide/1.0","User agent for the RETS client")
	userAgentPw := flag.String("user-agent-pw","","User agent authentication")
	logFile := flag.String("log-file","","")

	flag.Parse()

	var logger io.WriteCloser = nil
	if *logFile != "" {
		logger,err := os.Create(*logFile)
		if err != nil {
			panic(err)
		}
		defer logger.Close()
	}
	// should we throw an err here too?
	session, err := gorets.NewSession(*username, *password, *userAgent, *userAgentPw, logger)
	if err != nil {
		panic(err)
	}

	capability, err := session.Login(*loginUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Login: ", capability.Login)
	fmt.Println("Metadata: ", capability.GetMetadata)
	fmt.Println("Search: ", capability.Search)
	fmt.Println("GetObject: ", capability.GetObject)

	err = session.Get(capability.Get)
	if err != nil {
		panic(err)
	}

	mUrl := capability.GetMetadata
	format := "COMPACT"
	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-SYSTEM", "0"})
//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-RESOURCE", "0"})
//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-CLASS", "ActiveAgent"})
//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-TABLE", "ActiveAgent:ActiveAgent"})

	req := gorets.SearchRequest{
		Url: capability.Search,
		Query: "((LocaleListingStatus=|ACTIVE-CORE))",
		SearchType: "Property",
		Class: "ALL",
		Format: "COMPACT-DECODED",
		QueryType: "DMQL2",
		Count: gorets.COUNT_AFTER,
		Limit: 3,
		Offset: -1,
	}
	result, err := session.Search(req)
	if err != nil {
		panic(err)
	}
	cols := []string{
		"ListingKey",
		"ListPrice",
		"ListingID",
		"TotalPhotos",
		"ModificationTimestamp",
	}
	fmt.Println("COLUMNS:", cols)
	filter := result.FilterTo(cols)
	for row := range result.Data {
		fmt.Println(filter(row))
	}

	one,err := session.GetObject(gorets.GetObjectRequest{
		Url: capability.GetObject,
		Resource: "Property",
		Type: "Thumbnail",
		Id: "10385491290",
		ObjectId: "*",
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
	all,err := session.GetObject(gorets.GetObjectRequest{
		Url: capability.GetObject,
		Resource: "Property",
		Type: "Thumbnail",
		Id: "10388845716",
		ObjectId: "*",
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

	session.Logout(capability.Logout)
}
```
