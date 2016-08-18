package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jpfielding/gorets/rets"
	"golang.org/x/net/context"
)

func main2() {
	resource := flag.String("property", "Property", "Resource type for this object")
	id := flag.String("id", "1:*", "Comma separate list of ids: '234:*,123:0,123:1'")
	otype := flag.String("type", "Photo", "Object type for this request")
	directory := flag.String("directory", "", "Directory for file output")
	configFile := flag.String("config-file", "", "Config file for RETS connection")

	flag.Parse()

	config := Config{}
	if *configFile != "" {
		err := config.LoadFrom(*configFile)
		if err != nil {
			panic(err)
		}
	} else {
		// setup flag parsing for this type
		config.SetFlags()
		// flag.Parse()
	}

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		panic(err)
	}
	// setup timeout control
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// login
	capability, err := rets.Login(session, ctx, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	// make sure we close the rets connection
	defer rets.Logout(session, ctx, rets.LogoutRequest{URL: capability.Logout})
	// feedback
	fmt.Println("GetObject: ", capability.GetObject)
	// warning, this does _all_ of the photos
	one, err := rets.GetObjects(session, ctx, rets.GetObjectRequest{
		URL:      capability.GetObject,
		Resource: *resource,
		Type:     *otype,
		ID:       *id,
	})
	if err != nil {
		panic(err)
	}
	// TODO need to close chan, or better yet, not let it be a chan
	for r := range one {
		if err != nil {
			panic(err)
		}
		o := r.Object
		fmt.Println("PHOTO-META: ", o.ContentType, o.ContentID, o.ObjectID, len(o.Blob))
		// if we arent saving, then we quit
		if *directory == "" {
			continue
		}
		path := fmt.Sprintf("%s/%s", directory, o.ContentID)
		os.MkdirAll(path, os.ModePerm)
		f, err := os.Create(fmt.Sprintf("%s/%s", path, o.ObjectID))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.Write(o.Blob)
		if err != nil {
			panic(err)
		}
	}
}
