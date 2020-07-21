package main

import (
	encjson "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/rpc"
	"github.com/jpfielding/gorets/pkg/explorer"
	"github.com/rs/cors"
)

func main() {
	port := flag.String("port", "8000", "http port")
	cfgPath := flag.String("config", "", "config path")
	react := flag.String("react", "../../explorer/client/build", "ReactJS path")

	flag.Parse()

	cfg := config{Headers: []string{}}
	cfg.load(*cfgPath)

	// setup our mux for handling different paths
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*react)))

	// gorilla rpc endpoint services
	s := rpc.NewServer()
	s.RegisterService(&explorer.MetadataService{}, "")
	s.RegisterService(&explorer.SearchService{}, "")
	s.RegisterService(&explorer.ObjectService{}, "")
	corsOpts := cors.Options{
		AllowedOrigins: cfg.Headers,
		AllowedHeaders: []string{"OPTIONS", "POST", "GET", "HEAD"},
		ExposedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"},
	}
	// cors support wrapping our compressed(/rpc) services
	mux.Handle("/rpc", cors.New(corsOpts).Handler(handlers.CompressHandler(s)))

	// server startup
	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}

type config struct {
	Headers []string `json:"headers"`
}

func (c *config) load(path string) {
	if path == "" {
		return
	}
	tmp, err := importFrom(path)
	if err != nil {
		log.Fatal("Could not import Headers.")
	}
	c.Headers = tmp.Headers
}

func importFrom(path string) (config, error) {
	var cfgs config
	fmt.Println("walking:", path)
	walk := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(filename)
		if !strings.HasSuffix(base, ".json") {
			return nil
		}
		// eat all these errors since we're searching for matching json
		file, _ := os.Open(filename)
		blob, _ := ioutil.ReadAll(file)
		_ = file.Close()
		_ = encjson.Unmarshal(blob, &cfgs)
		return nil
	}
	err := filepath.Walk(path, walk)
	return cfgs, err
}
