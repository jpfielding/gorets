package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"context"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
)

func main() {
	status := 0
	defer func() { os.Exit(status) }()

	optionsFile := flag.String("object-options", "", "Get object")
	configFile := flag.String("config-file", "", "Config file for RETS connection")
	output := flag.String("output", "", "Directory for file output")

	config := common.Config{}
	config.SetFlags()

	getOptions := GetObjectOptions{}
	getOptions.SetFlags()

	flag.Parse()

	if *configFile != "" {
		err := config.LoadFrom(*configFile)
		if err != nil {
			log.Println(err)
			status = 1
			return
		}
	}
	log.Printf("Connection Settings: %v\n", config)
	if *optionsFile != "" {
		err := getOptions.LoadFrom(*optionsFile)
		if err != nil {
			log.Println(err)
			status = 2
			return
		}
	}
	log.Printf("GetObject Options: %v\n", getOptions)

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		log.Println(err)
		status = 3
		return
	}
	// setup timeout control
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// login
	capability, err := rets.Login(session, ctx, rets.LoginRequest{URL: config.URL})
	if err != nil {
		log.Println(err)
		status = 4
		return
	}
	// make sure we close the rets connection
	defer rets.Logout(session, ctx, rets.LogoutRequest{URL: capability.Logout})
	// feedback
	log.Println("GetObject: ", capability.GetObject)
	// warning, this does _all_ of the photos
	response, err := rets.GetObjects(session, ctx, rets.GetObjectRequest{
		URL: capability.GetObject,
		GetObjectParams: rets.GetObjectParams{
			Resource: getOptions.Resource,
			Type:     getOptions.Type,
			ID:       getOptions.ID,
		},
	})
	if err != nil {
		log.Println(err)
		status = 5
		return
	}
	defer response.Close()
	err = response.ForEach(func(o *rets.Object, err error) error {
		log.Println("PHOTO-META: ", o.ContentType, o.ContentID, o.ObjectID, len(o.Blob))
		// if we arent saving, then we quit
		if *output == "" {
			return nil
		}
		path := fmt.Sprintf("%s/%s", *output, o.ContentID)
		os.MkdirAll(path, os.ModePerm)
		f, err := os.Create(fmt.Sprintf("%s/%d", path, o.ObjectID))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(o.Blob)
		return err
	})
	if err != nil {
		log.Println(err)
		status = 6
		return
	}
}

// GetObjectOptions ...
type GetObjectOptions struct {
	Resource string `json:"resource"`
	Type     string `json:"type"`
	ID       string `json:"id"`
}

// SetFlags ...
func (o *GetObjectOptions) SetFlags() {
	flag.StringVar(&o.Resource, "resource", "Property", "Resource for the search")
	flag.StringVar(&o.Type, "type", "Photo", "Photo, document, etc...")
	flag.StringVar(&o.ID, "id", "*", "Subtype of resource")
}

// LoadFrom ...
func (o *GetObjectOptions) LoadFrom(filename string) error {
	// xlog.Println("loading:", filename)
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, o)
	if err != nil {
		return err
	}
	return nil
}
