package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jpfielding/gorets/pkg/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"github.com/spf13/cobra"
)

const (
	connectFlag  = "connect"
	connectUsage = "RETS Configuration. If not specified, assumes $HOME/.gorets/connect.json"
	outputFlag   = "output"
	outputUsage  = "Output directory"
	timeoutFlag  = "tiemout"
	timeoutUsage = "Request timeout (in seconds)"
	wirelogFlag  = "wirelog"
	wirelogUsage = "The file that will contain wirelog output"
)

func main() {
	NewCmd().Execute()
}

// NewCmd ...
func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gorets [command]",
		Short: "gorets is used to communicate with RETS servers",
	}
	cmd.AddCommand(
		NewMetadataCmd(),
		NewSearchCmd(),
		NewGetObjectsCmd(),
	)
	cmd.PersistentFlags().StringP(connectFlag, "c", "config.json", connectUsage)
	cmd.PersistentFlags().StringP(wirelogFlag, "w", "wire.log", wirelogUsage)
	cmd.PersistentFlags().StringP(outputFlag, "o", "", outputUsage)
	cmd.PersistentFlags().Int64P(timeoutFlag, "t", 60, timeoutUsage)
	return cmd
}

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
	if err != nil {
		return err
	}
	err = json.Unmarshal(blob, model)
	if err != nil {
		return err
	}
	return nil
}

// getPersistentFlagValues extracts the persistent flag values for us in each command
func getPersistentFlagValues(f interface{}, cmd *cobra.Command) (Connect, string, time.Duration) {
	cFile, err := cmd.Flags().GetString(connectFlag)
	handleError(f, err)

	output, err := cmd.Flags().GetString(outputFlag)
	handleError(f, err)

	// TODO investigate using GetDuration
	timeout, err := cmd.Flags().GetInt64(timeoutFlag)
	handleError(f, err)

	// the params for talking to the server
	connect := Connect{}
	err = LoadFrom(cFile, &connect)
	handleError(f, err)

	return connect, output, (time.Duration(timeout) * time.Second)
}

func fmtError(function interface{}, err error) error {
	if err == nil {
		return nil
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	splitName := strings.Split(fullName, ".")
	localName := splitName[len(splitName)-1]
	return fmt.Errorf("%s: %v", localName, err)
}

func handleError(function interface{}, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", fmtError(function, err))
		os.Exit(1)
	}
}
