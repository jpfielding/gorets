package main

import (
	"context"
	"io"
	"os"

	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

const (
	mTypeFlag   = "type"
	mFormatFlag = "format"
	mIDFlag     = "id"
)

// NewMetadataCmd ...
func NewMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Extract RETS metadata",
		Run:   metadata,
	}
	cmd.Flags().String(mTypeFlag, "METADATA-SYSTEM", "The type of metadata requested")
	cmd.Flags().String(mFormatFlag, "COMPACT", "Metadata format")
	cmd.Flags().String(mIDFlag, "*", "Metadata identifier")
	return cmd
}

func metadata(cmd *cobra.Command, args []string) {
	connect, output, timeout := getPersistentFlagValues(search, cmd)

	var err error // annoying

	params := MetadataOptions{}
	params.MType, err = cmd.Flags().GetString(mTypeFlag)
	handleError(metadata, err)

	params.Format, err = cmd.Flags().GetString(mFormatFlag)
	handleError(metadata, err)

	params.ID, err = cmd.Flags().GetString(mIDFlag)
	handleError(metadata, err)

	// should we throw an err here too?
	session, err := connect.Initialize()
	handleError(metadata, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	capability, err := rets.Login(ctx, session, rets.LoginRequest{URL: connect.URL})
	handleError(metadata, err)

	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: capability.Logout})
	mp := rets.MetadataParams{
		Format: params.Format,
		MType:  params.MType,
		ID:     params.ID,
	}
	reader, err := rets.MetadataStream(rets.MetadataResponse(ctx, session, rets.MetadataRequest{
		URL:            capability.GetMetadata,
		MetadataParams: mp,
	}))
	defer reader.Close()
	handleError(metadata, err)

	out := os.Stdout
	if output != "" {
		out, err = os.Create(output + "/metadata.xml")
		handleError(metadata, err)
		defer out.Close()
	}
	io.Copy(out, reader)
}

// MetadataOptions ...
type MetadataOptions struct {
	MType  string `json:"metadata-type"`
	Format string `json:"format"`
	ID     string `json:"id"`
}
