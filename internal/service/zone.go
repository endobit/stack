package service

import (
	"context"
	"fmt"

	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
)

// ListZones returns a list of zones.
func (s *service) ListZones(r *pb.ListZonesRequest, in pb.StackService_ListZonesServer) error {
	zones, err := s.Queries.GetZones(context.Background(), r.ZoneGlob)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	for _, z := range zones {
		err := in.Send(&pb.ListZonesResponse{
			Id:       z.ID,
			Name:     z.Zone,
			TimeZone: z.TimeZone,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
