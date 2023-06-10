package service

import (
	"net"
	"os"
	"strconv"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"

	"github.com/endobit/clog"
	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
	"github.com/endobit/stack/internal/service/store"
)

// NewRootCmd returns a new root command.
func NewRootCmd(version string) *cobra.Command {
	var port int

	cmd := cobra.Command{
		Use:     "stackd",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				return err
			}

			queries, err := store.Connect()
			if err != nil {
				return err
			}

			svr := grpc.NewServer()
			svc := service{
				Logger:  slog.New(clog.NewHandler(os.Stderr)),
				Queries: queries,
			}

			pb.RegisterStackServiceServer(svr, &svc)
			return svr.Serve(listener)
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")

	return &cmd
}
