package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// Connect info to use the demo
type Connect struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	URL         string `json:"url"`
	UserAgent   string `json:"user-agent"`
	UserAgentPw string `json:"user-agent-pw"`
	Version     string `json:"rets-version"`
	WireLog     string `json:"wirelog"`
}

// Initialize extracts the cmd line params and creates the rets.Requester
func (cnt *Connect) Initialize() (rets.Requester, error) {
	transport := wirelog.NewHTTPTransport()

	if cnt.WireLog != "" {
		wirelog.LogToFile(transport, cnt.WireLog, true, true)
		fmt.Println("wire logging enabled:", cnt.WireLog)
	}
	// should we throw an err here too?
	return rets.DefaultSession(
		cnt.Username,
		cnt.Password,
		cnt.UserAgent,
		cnt.UserAgentPw,
		cnt.Version,
		transport,
	)
}

// LoadFrom loads the model onto the struct
func LoadFrom(filename string, model interface{}) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, model)
	if err != nil {
		return err
	}
	return nil
}
