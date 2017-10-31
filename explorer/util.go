package explorer

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/jpfielding/gorets/rets"
	"github.com/jpfielding/gowirelog/wirelog"
	"golang.org/x/net/proxy"
)

// Config ...
type Config struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserAgent   string `json:"userAgent"`
	UserAgentPw string `json:"userAgentPw"`
	Proxy       string `json:"proxy"`
	Version     string `json:"version"`
}

// ReadWirelog ...
func (c Config) ReadWirelog(fun func(*os.File, error) error) error {
	f, err := os.Open(c.Wirelog())
	defer f.Close()
	if e := fun(f, err); e != nil {
		return e
	}
	return nil
}

// TODO creating a unqiue tmp path and storing this instantiation there

// Wirelog path
func (c Config) Wirelog() string {
	return fmt.Sprintf("/tmp/gorets/%s/wire.log", c.ID)
}

// MSystem path
func (c Config) MSystem() string {
	return fmt.Sprintf("/tmp/gorets/%s/metadata.json", c.ID)
}

// Connect ...
// TODO need to remove connecting from session creation
func (c Config) Connect(ctx context.Context) (*Session, error) {
	// start with the default Dialer from http.DefaultTransport
	transport := wirelog.NewHTTPTransport()
	// if there is a need to proxy
	if c.Proxy != "" {
		log.Printf("Using proxy %s", c.Proxy)
		d, err := proxy.SOCKS5("tcp", c.Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("rets proxy: %v", err)
		}
		transport.Dial = d.Dial
	}
	var closer io.Closer
	var err error
	// wire logging
	logFile := c.Wirelog()
	closer, err = wirelog.LogToFile(transport, logFile, true, true)
	if err != nil {
		return nil, fmt.Errorf("wirelog setup: %v", err)
	}
	log.Printf("wire logging enabled %s", logFile)

	// should we throw an err here too?
	sess, err := rets.DefaultSession(
		c.Username,
		c.Password,
		c.UserAgent,
		c.UserAgentPw,
		c.Version,
		transport,
	)
	urls, err := rets.Login(ctx, sess, rets.LoginRequest{URL: c.URL})
	if err != nil {
		return nil, err
	}
	return &Session{
		requester: sess,
		closer:    closer,
		urls:      *urls,
	}, err
}

// Session ...
type Session struct {
	Config Config // TODO should this be private?
	// things in user state
	closer    io.Closer
	requester rets.Requester
	urls      rets.CapabilityURLs
}

// Close is an io.Closer
func (s *Session) Close(ctx context.Context) error {
	_, err := rets.Logout(ctx, s.requester, rets.LogoutRequest{URL: s.urls.Logout})
	s.closer = nil
	s.requester = nil
	s.closer.Close()
	return err
}

// Op is simple Rets function
type Op func(rets.Requester, rets.CapabilityURLs) error

// Process processes a set of requests
// TODO needs a retry mechanism
// TOOD deal with keeping sessions alive
func (s *Session) Process(ctx context.Context, ops ...Op) error {
	for _, op := range ops {
		err := op(s.requester, s.urls)
		if err != nil {
			return err
		}
	}
	return nil
}

// GZIP the output

//JSONExist ...
func JSONExist(filename string) bool {
	_, err := os.Stat(filename + ".gz")
	return !os.IsNotExist(err)
}

// JSONStore raw file storage
func JSONStore(filename string, data interface{}) error {
	dir := path.Dir(filename)
	// TODO dont repeat this for every write
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// formatted
	var out bytes.Buffer
	json.Indent(&out, raw, "", "\t")
	z := gzip.NewWriter(f)
	_, err = out.WriteTo(z)
	if err != nil {
		return err
	}
	err = z.Close()
	if err != nil {
		return err
	}
	f.Close()
	err = os.Rename(f.Name(), filename+".gz")
	if err != nil {
		return err
	}
	fmt.Println("wrote:", filename)
	return nil
}

// JSONLoad raw file load
func JSONLoad(filename string, data interface{}) error {
	file, err := os.Open(filename + ".gz")
	defer file.Close()
	if err != nil {
		return err
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(gz)
	err = json.Unmarshal(blob, data)
	if err != nil {
		return err
	}
	return nil
}
