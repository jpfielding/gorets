package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jpfielding/gorets/proxy"
)

// TOOD read in login params, get source config from url path
// TODO load gorets to fulfull request, proxy response back
// TODO cache source sessions by auth info
// TODO store caching options (concurrent or login limits) per config

func main() {
	bind := flag.String("bind", ":8000", "The host:port to bind this service")
	config := flag.String("config", "config.json", "The configuration for this service")

	flag.Parse()

	// load from config
	var cfgs []proxy.Config
	raw, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &cfgs)
	if err != nil {
		panic(err)
	}

	prefix := "/proxy/rets"
	ops := map[string]string{
		"Login":       fmt.Sprintf("%s/login/", prefix),
		"GetMetadata": fmt.Sprintf("%s/metadata/", prefix),
		"Search":      fmt.Sprintf("%s/search/", prefix),
		"GetObject":   fmt.Sprintf("%s/getobject/", prefix),
	}
	srcs := proxy.NewSources(cfgs)
	// the base /rets/ proxy handler
	http.HandleFunc(ops["Login"], proxy.Login(ops, srcs))
	http.HandleFunc(ops["GetMetadata"], proxy.Metadata(ops, srcs))
	http.HandleFunc(ops["Search"], proxy.Search(ops, srcs))
	// web server
	err = http.ListenAndServe(*bind, nil)
	if err != nil {
		panic(err)
	}
}
