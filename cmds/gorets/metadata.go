package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

func NewMetadataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search a RETS server",
		Run:   metadata,
	}
}

func metadata(cmd *cobra.Command, args []string) {
	config := cmd.Flags().GetString("config", "", "Path to the config info for RETS connection")
	output := cmd.Flags().GetString("output", "", "Directory for file output")
	timeout := cmd.Flags().GetInt("timeout", 60, "Seconds to timeout the connection")

	connect := common.Connect{}
	LoadFrom(connectFile, &connect)

	metadataOpts := MetadataOptions{}
	LoadFrom(metadataFile, &metadataOpts)
	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Seconds)
	defer cancel()

	capability, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: capability.Logout})

	reader, err := rets.MetadataStream(ctx, session, rets.MetadataRequest{
		URL:    capability.GetMetadata,
		Format: metadataOpts.Format,
		MType:  metadataOpts.MType,
		ID:     metadataOpts.ID,
	})
	defer reader.Close()
	if err != nil {
		panic(err)
	}
	out := os.Stdout
	if *output != "" {
		out, _ = os.Create(*output + "/metadata.xml")
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

func (o *MetadataOptions) GetFlags() {
	o.MType = cmd.Flags().GetString("mtype", "METADATA-SYSTEM", "The type of metadata requested")
	o.Format = cmd.Flags().GetString("format", "COMPACT", "Metadata format")
	o.ID = cmd.Flags().GetString("id", "*", "Metadata identifier")
}
