package common

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
)

// Config info to use the demo
type Config struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	URL         string `json:"url"`
	UserAgent   string `json:"user-agent"`
	UserAgentPw string `json:"user-agent-pw"`
	Version     string `json:"rets-version"`
	WireLog     string `json:"wirelog"`
}

// Initialize extracts the cmd line params and creates the rets.Requester
func (cfg *Config) Initialize() (rets.Requester, error) {
	transport := wirelog.NewHTTPTransport()

	if cfg.WireLog != "" {
		wirelog.LogToFile(transport, cfg.WireLog, true, true)
		log.Println("wire logging enabled:", cfg.WireLog)
	}
	// should we throw an err here too?
	return rets.DefaultSession(
		cfg.Username,
		cfg.Password,
		cfg.UserAgent,
		cfg.UserAgentPw,
		cfg.Version,
		transport,
	)
}

// LoadFrom ...
func (cfg *Config) LoadFrom(filename string) error {
	// xlog.Println("loading:", filename)
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, cfg)
	if err != nil {
		return err
	}
	return nil
}

// SetFlags bings this Config to flags from the command line
func (cfg *Config) SetFlags() {
	flag.StringVar(&cfg.Username, "username", "", "Username for the RETS server")
	flag.StringVar(&cfg.Password, "password", "", "Password for the RETS server")
	flag.StringVar(&cfg.URL, "Config-url", "", "Config URL for the RETS server")
	flag.StringVar(&cfg.UserAgent, "user-agent", "Threewide/1.0", "User agent for the RETS client")
	flag.StringVar(&cfg.UserAgentPw, "user-agent-pw", "", "User agent authentication")
	flag.StringVar(&cfg.Version, "rets-version", "", "RETS Version")
	flag.StringVar(&cfg.WireLog, "log-file", "", "")
}
