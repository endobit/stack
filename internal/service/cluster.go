package service

import (
	"context"
	"fmt"

	"github.com/endobit/stack/gen/go/db"
	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
)

// ListClusters lists all the clusters.
func (s *service) ListClusters(r *pb.ListClustersRequest, in pb.StackService_ListClustersServer) error {
	clusters, err := s.Queries.GetClusters(context.Background(),
		db.GetClustersParams{
			Zone:    r.ZoneGlob,
			Cluster: r.ClusterGlob,
		})
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	for _, c := range clusters {
		err := in.Send(&pb.ListClustersResponse{
			Id:   c.ID,
			Name: c.Cluster,
			Zone: c.Zone,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
