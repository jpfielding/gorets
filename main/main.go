/**
 * Created with IntelliJ IDEA.
 * User: jp
 * Date: 9/25/13
 * Time: 10:02 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	gorets "github.com/jpfielding/gorets_client"
)

func main() {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	loginUrl := flag.String("login-url", "", "Login URL for the RETS server")
	userAgent := flag.String("user-agent", "Threewide/1.0", "User agent for the RETS client")
	userAgentPw := flag.String("user-agent-pw", "", "User agent authentication")
	retsVersion := flag.String("rets-version", "", "RETS Version")
	logFile := flag.String("log-file", "", "")

	flag.Parse()

	var logger io.WriteCloser = nil
	if *logFile != "" {
		file, err := os.Create(*logFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logger = file
		fmt.Println("wire logging enabled: ", file.Name())
	}
	// should we throw an err here too?
	session, err := gorets.NewSession(*username, *password, *userAgent, *userAgentPw, *retsVersion, logger)
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
		fmt.Println("this was stupid, shouldnt even be here")
	}

	mUrl := capability.GetMetadata
	format := "COMPACT"
	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-SYSTEM", "0"})
	//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-RESOURCE", "0"})
	//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-CLASS", "ActiveAgent"})
	//	session.GetMetadata(gorets.MetadataRequest{mUrl, format, "METADATA-TABLE", "ActiveAgent:ActiveAgent"})

	req := gorets.SearchRequest{
		Url:        capability.Search,
		Query:      "((180=|AH))",
		SearchType: "Property",
		Class:      "1",
		Format:     "COMPACT-DECODED",
		QueryType:  "DMQL2",
		Count:      gorets.COUNT_AFTER,
		Limit:      3,
		Offset:     -1,
	}
	result, err := session.Search(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("COLUMNS:", result.Columns)
	for row := range result.Data {
		fmt.Println(row)
	}

	one, err := session.GetObject(gorets.GetObjectRequest{
		Url:      capability.GetObject,
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
	all, err := session.GetObject(gorets.GetObjectRequest{
		Url:      capability.GetObject,
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

	session.Logout(capability.Logout)
}
