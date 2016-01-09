package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/jpfielding/gorets/client"
)

func main() {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	format := flag.String("format", "COMPACT", "COMPACT or STANDARD-XML")
	mType := flag.String("type", "METADATA-SYSTEM", "Metadata Type 'METADATA-TABLE'")
	id := flag.String("id", "*", "Property, Office, Agent, User")
	loginURL := flag.String("login-url", "http://sum.rets.interealty.com/Login.asmx/Login", "Login URL for the RETS server")
	userAgent := flag.String("user-agent", "Threewide/1.5", "User agent for the RETS client")
	userAgentPw := flag.String("user-agent-pw", "listhub", "User agent authentication")
	retsVersion := flag.String("rets-version", "RETS/1.5", "RETS Version")
	logFile := flag.String("log-file", "/tmp/listhub/rets.log", "")

	flag.Parse()

	d := net.Dial

	if *logFile != "" {
		file, err := os.Create(*logFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		fmt.Println("wire logging enabled: ", file.Name())
		d = client.WireLog(file, d)
	}

	transport := http.Transport{
		DisableCompression: true,
		Dial:               d,
	}

	// should we throw an err here too?
	session, err := client.DefaultSession(
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

	capability, err := client.Login(session, ctx, client.LoginRequest{URL: *loginURL})
	if err != nil {
		panic(err)
	}
	fmt.Println("Login: ", capability.Login)
	fmt.Println("Metadata: ", capability.GetMetadata)

	mURL := capability.GetMetadata

	metadata, err := client.GetCompactMetadata(session, ctx, client.MetadataRequest{
		URL:    mURL,
		Format: *format,
		MType:  *mType,
		ID:     *id,
	})

	if err != nil {
		panic(err)
	}

	for id, elems := range metadata.Elements {
		for _, cd := range elems {
			fmt.Println(id + ": " + cd.ID)
		}
	}
}
