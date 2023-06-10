// Package service implements the stack grpc service.
package service

import (
	"golang.org/x/exp/slog"

	"github.com/endobit/stack/gen/go/db"
	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
)

type service struct {
	pb.UnimplementedStackServiceServer
	Logger  *slog.Logger
	Queries *db.Queries
}
