package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
	"golang.org/x/net/context"
)

func main() {
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
			panic(err)
		}
	}
	fmt.Printf("Connection Settings: %v\n", config)
	if *optionsFile != "" {
		err := getOptions.LoadFrom(*optionsFile)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("GetObject Options: %v\n", getOptions)

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
		Resource: getOptions.Resource,
		Type:     getOptions.Type,
		ID:       getOptions.ID,
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
		if *output == "" {
			continue
		}
		path := fmt.Sprintf("%s/%d", *output, o.ContentID)
		os.MkdirAll(path, os.ModePerm)
		f, err := os.Create(fmt.Sprintf("%s/%d", path, o.ObjectID))
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
