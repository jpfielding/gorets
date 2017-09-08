package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"context"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

func NewGetObjectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "getobjects",
		Short: "Get Objects from a RETS server",
		Run:   getObjects,
	}
}

func getObjects(cmd *cobra.Command, args []string) {
	config := cmd.Flags().GetString("config", "", "Path to the config info for RETS connection")
	output := cmd.Flags().GetString("output", "", "Directory for file output")
	timeout := cmd.Flags().GetInt("timeout", 60, "Seconds to timeout the connection")

	connect := common.Connect{}
	LoadFrom(connectFile, &connect)

	getOptions := GetObjectOptions{}
	LoadFrom(optionsFile, &getOptions)

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		panic(err)
	}
	// setup timeout control
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Seconds)
	defer cancel()
	// login
	capability, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	// make sure we close the rets connection
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: capability.Logout})
	// feedback
	fmt.Println("GetObject: ", capability.GetObject)
	// warning, this does _all_ of the photos
	response, err := rets.GetObjects(ctx, session, rets.GetObjectRequest{
		URL: capability.GetObject,
		GetObjectParams: rets.GetObjectParams{
			Resource: getOptions.Resource,
			Type:     getOptions.Type,
			ID:       getOptions.ID,
		},
	})
	if err != nil {
		panic(err)
	}
	defer response.Close()
	err = response.ForEach(func(o *rets.Object, err error) error {
		fmt.Println("PHOTO-META: ", o.ContentType, o.ContentID, o.ObjectID, len(o.Blob))
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
		panic(err)
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
	blob, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, o)
	if err != nil {
		return err
	}
	return nil
}
