package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/jpfielding/gorets/client"
)

func main() {
	username := flag.String("username", "", "Username for the RETS server")
	password := flag.String("password", "", "Password for the RETS server")
	loginURL := flag.String("login-url", "http://sum.rets.interealty.com/Login.asmx/Login", "Login URL for the RETS server")
	userAgent := flag.String("user-agent", "Threewide/1.5", "User agent for the RETS client")
	userAgentPw := flag.String("user-agent-pw", "listhub", "User agent authentication")
	retsVersion := flag.String("rets-version", "RETS/1.5", "RETS Version")
	logFile := flag.String("log-file", "/tmp/listhub/sarco-rets.log", "")

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

	// should we throw an err here too?
	session, err := client.NewSession(*username, *password, *userAgent, *userAgentPw, *retsVersion, &http.Transport{
		DisableCompression: true,
		Dial:               d,
	})
	if err != nil {
		panic(err)
	}

	capability, err := session.Login(client.LoginRequest{URL: *loginURL})
	if err != nil {
		panic(err)
	}
	fmt.Println("Login: ", capability.Login)
	fmt.Println("Metadata: ", capability.GetMetadata)
	fmt.Println("Search: ", capability.Search)
	fmt.Println("GetObject: ", capability.GetObject)

	err = session.Get(client.GetRequest{URL: capability.Get})
	if err != nil {
		fmt.Println("this was stupid, shouldnt even be here")
	}

	mURL := capability.GetMetadata

	for _, f := range []string{"STANDARD-XML", "COMPACT"} {
		for _, t := range []string{"TABLE"} {
			session.GetMetadata(client.MetadataRequest{
				URL:    mURL,
				Format: f,
				MType:  "METADATA-" + t,
				Id:     "Office",
			})
		}
	}
}
