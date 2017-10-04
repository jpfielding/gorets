package main

import (
	"fmt"
	"os"

	"context"

	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

const (
	goResourceFlag = "resource"
	goTypeFlag     = "type"
	goIDFlag       = "id"
)

// NewGetObjectsCmd ...
func NewGetObjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getobjects",
		Short: "Get Objects from a RETS server",
		Run:   getObjects,
	}
	cmd.Flags().String(goResourceFlag, "Property", "Resource type of the request")
	cmd.Flags().String(goTypeFlag, "Photo", "Metadata format")
	cmd.Flags().String(goIDFlag, "", "Object ID for the request ('key1:*;key2:0')")

	return cmd
}

func getObjects(cmd *cobra.Command, args []string) {
	connect, output, timeout := getPersistentFlagValues(search, cmd)

	var err error // annoying

	params := GetObjectParams{}

	params.Resource, err = cmd.Flags().GetString(goResourceFlag)
	handleError(getObjects, err)

	params.Type, err = cmd.Flags().GetString(goTypeFlag)
	handleError(getObjects, err)

	params.ID, err = cmd.Flags().GetString(goIDFlag)
	handleError(getObjects, err)

	// should we throw an err here too?
	session, err := connect.Initialize()
	handleError(getObjects, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	capability, err := rets.Login(ctx, session, rets.LoginRequest{URL: connect.URL})
	handleError(getObjects, err)
	// make sure we close the rets connection
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: capability.Logout})

	// feedback
	fmt.Println("GetObject: ", capability.GetObject)
	// warning, this does _all_ of the photos
	gor, err := rets.GetObjects(ctx, session, rets.GetObjectRequest{
		URL: capability.GetObject,
		GetObjectParams: rets.GetObjectParams{
			Resource: params.Resource,
			Type:     params.Type,
			ID:       params.ID,
		},
	})
	handleError(getObjects, err)
	response := &rets.GetObjectResponse{Response: gor}
	defer response.Close()
	err = response.ForEach(func(o *rets.Object, err error) error {
		fmt.Println("PHOTO-META: ", o.ContentType, o.ContentID, o.ObjectID, len(o.Blob))
		// if we arent saving, then we quit
		if output == "" {
			return nil
		}
		path := fmt.Sprintf("%s/%s", output, o.ContentID)
		os.MkdirAll(path, os.ModePerm)
		f, err := os.Create(fmt.Sprintf("%s/%d", path, o.ObjectID))
		handleError(getObjects, err)

		defer f.Close()
		_, err = f.Write(o.Blob)
		return err
	})
	handleError(getObjects, err)
}

// GetObjectParams ...
type GetObjectParams struct {
	Resource string `json:"resource"`
	Type     string `json:"type"`
	ID       string `json:"id"`
}
