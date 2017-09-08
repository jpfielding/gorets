package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jpfielding/gorets/rets"
	"github.com/spf13/cobra"
)

func NewGetPayloadListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "getpayloadlist",
		Short: "Get a payload list from a RETS server",
		Run:   getPayloadList,
	}
}

func getPayloadList(cmd *cobra.Command, args []string) {
	config := cmd.Flags().GetString("config", "", "Path to the config info for RETS connection")
	output := cmd.Flags().GetString("output", "", "Directory for file output")
	timeout := cmd.Flags().GetInt("timeout", 60, "Seconds to timeout the connection")

	// the params for talking to the server
	connect := Connect{}
	LoadFrom(connectFile, &connect)

	// launch the session
	session, err := connect.Initialize()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Seconds)
	defer cancel()

	urls, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		panic(err)
	}
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: urls.Logout})

	if urls.GetPayloadList != "" {

	}
	fmt.Println("Payloads: ", urls.GetPayloadList)
	payloads, err := rets.GetPayloadList(ctx, session, rets.PayloadListRequest{
		URL: urls.GetPayloadList,
		ID:  fmt.Sprintf("%s:%s", searchOpts.Resource, searchOpts.Class),
	})
	if err != nil {
		panic(err)
	}
	err = payloads.ForEach(func(payload rets.CompactData, err error) error {
		fmt.Printf("%v\n", payload)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func processCompact(ctx context.Context, sess rets.Requester, req rets.SearchRequest, output *string) {
	w := csv.NewWriter(os.Stdout)
	if *output != "" {
		os.MkdirAll(*output, 0777)
		f, _ := os.Create(*output + "/results.csv")
		defer f.Close()
		w = csv.NewWriter(f)
	}
	defer w.Flush()

	// loop over all the pages we need
	for {
		fmt.Printf("Querying next page: %v\n", req)
		result, err := rets.SearchCompact(ctx, sess, req)
		if err != nil {
			panic(err)
		}
		switch result.Response.Code {
		case rets.StatusOK:
			// we got some daters
		case rets.StatusNoRecords:
			return
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			panic(errors.New(result.Response.Text))
		}
		count := 0
		if count == 0 {
			w.Write(result.Columns)
		}
		hasMoreRows, err := result.ForEach(func(row rets.Row, err error) error {
			if err != nil {
				return err
			}
			w.Write(row)
			count++
			return err
		})
		result.Close()
		if err != nil {
			panic(err)
		}
		if !hasMoreRows {
			return
		}
		if req.Offset == 0 {
			req.Offset = 1
		}
		req.Offset = req.Offset + count
	}
}
