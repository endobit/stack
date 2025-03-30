// Package main implements the stack CLI.
package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"endobit.io/metal"
	"endobit.io/metal/logging"
	"endobit.io/mops"
	"endobit.io/stack/internal/commands"
)

var version string

func main() {
	cmd := newRootCmd()
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		metalUser, metalPass, metalServer string
		metalClient                       metal.Client
		mopsServer                        string
		mopsClient                        mops.Client
		logOpts                           *logging.Options
	)

	cmd := cobra.Command{
		Use:   "stack",
		Short: "Stack Client",
		Long:  "Stack Command Line Client",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			logger, err := logOpts.NewLogger()
			if err != nil {
				return err
			}

			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				MinVersion:         tls.VersionTLS12,
			})

			conn, err := grpc.NewClient(metalServer, grpc.WithTransportCredentials(creds))
			if err != nil {
				return err
			}

			metalClient = *metal.NewClient(conn, logger)
			if err := metalClient.Authorize(metalUser, metalPass); err != nil {
				return err
			}

			mopsClient = mops.Client{
				URL: "http://" + mopsServer,
				Client: http.Client{
					Timeout: 5 * time.Second,
				},
			}

			return nil
		},
	}

	logOpts = logging.NewOptions(cmd.PersistentFlags())

	cmd.PersistentFlags().StringVar(&metalUser, "metal-user", "admin", "username for metal authentication")
	cmd.PersistentFlags().StringVar(&metalPass, "metal-pass", "admin", "password for metal authentication")
	cmd.PersistentFlags().StringVar(&metalServer, "metal", "localhost:"+strconv.Itoa(metal.DefaultPort),
		"address of the metal server")
	cmd.PersistentFlags().StringVar(&mopsServer, "mops", "localhost:"+strconv.Itoa(mops.DefaultPort),
		"address of the mops server")

	root := commands.Root{
		Metal: &metalClient,
		Ops:   &mopsClient,
	}

	cmd.AddCommand(
		root.New(commands.Add),
		root.New(commands.Dump),
		root.New(commands.List),
		root.New(commands.Load),
		root.New(commands.Remove),
		root.New(commands.Report),
		root.New(commands.Set),
		root.New(commands.Unset))

	return &cmd
}
