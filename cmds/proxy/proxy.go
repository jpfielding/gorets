package main

import (
	"encoding/json"
	"flag"
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
	cfgs := make(map[string][]proxy.Config)
	raw, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &cfgs)
	if err != nil {
		panic(err)
	}

	// the base /rets/ proxy handler
	login := "/proxy/rets/login/"
	http.HandleFunc(login, proxy.Login(login, cfgs))

	err = http.ListenAndServe(*bind, nil)
	if err != nil {
		panic(err)
	}
}
