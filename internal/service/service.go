// Package service implements the stack grpc service.
package service

import (
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/endobit/stack/gen/go/db"
	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
)

type service struct {
	pb.UnimplementedStackServiceServer
	Logger  *slog.Logger
	Queries *db.Queries
}

// QueryError is returned when a query fails.
type QueryError struct {
	Err error
}

// Error implements the error interface.
func (e QueryError) Error() string {
	return fmt.Sprintf("query failed: %v", e.Err)
}

// bool2int converts a bool to an int64.
func bool2int(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

// int2bool converts an int64 to a bool.
func int2bool(i int64) bool {
	return i == 1
}
