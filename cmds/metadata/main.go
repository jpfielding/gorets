package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
	"golang.org/x/net/context"
)

func main() {
	mtype := flag.String("mtype", "METADATA-SYSTEM", "The type of metadata requested")
	format := flag.String("format", "COMPACT", "Metadata format")
	id := flag.String("id", "*", "Metadata identifier")

	configFile := flag.String("config-file", "", "Config file for RETS connection")

	config := common.Config{}
	config.SetFlags()

	flag.Parse()

	if *configFile != "" {
		err := config.LoadFrom(*configFile)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Connection Settings: %v\n", config)

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	capability, err := rets.Login(session, ctx, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	defer rets.Logout(session, ctx, rets.LogoutRequest{URL: capability.Logout})

	reader, err := rets.MetadataStream(session, ctx, rets.MetadataRequest{
		URL:    capability.GetMetadata,
		Format: *format,
		MType:  *mtype,
		ID:     *id,
	})
	if err != nil {
		panic(err)
	}
	// TODO do something meaningful with this
	reader.Close()
}
