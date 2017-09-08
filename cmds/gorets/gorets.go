package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	NewCmd().Execute()
}

// NewCmd ...
func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gorets",
		Short: "gorets is used to communicate with RETS servers",
	}
	GoRETSCmd.AddCommand(
		NewMetadataCmd(),
		NewGetPayloadListCmd(),
		NewSearchCmd(),
		NewGetObjectCmd(),
	)
	// cmd.PersistentFlags().StringVarP(&ConfigURL, "server", "s", config.GetServiceURL(), "config service URL")
	return cmd
}

func dieOnError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
