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
	"github.com/gorilla/rpc/json"
	"github.com/jpfielding/gorets/explorer"
)

type config struct {
	Headers []string `json:"headers"`
}

func main() {
	port := flag.String("port", "8000", "http port")
	cnflag := flag.String("config", "", "config path")
	react := flag.String("react", "../../explorer/client/build", "ReactJS path")

	flag.Parse()

	cfg := config{Headers: []string{}}
	if *cnflag != "" {
		tmp, err := importFrom(*cnflag)
		if err != nil {
			log.Fatal("Could not import Headers.")
		} else {
			cfg = tmp
		}
	}

	http.Handle("/", http.FileServer(http.Dir(*react)))

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(explorer.CodecWithCors(cfg.Headers, json.NewCodec()), "application/json")
	s.RegisterService(&explorer.MetadataService{}, "")
	s.RegisterService(&explorer.SearchService{}, "")
	s.RegisterService(&explorer.ObjectService{}, "")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "HEAD"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"}),
	)
	// rpc calls
	http.Handle("/rpc", handlers.CompressHandler(cors(s)))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
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
