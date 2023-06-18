package service

import (
	"context"

	"github.com/endobit/stack/gen/go/db"
	pb "github.com/endobit/stack/gen/go/proto/stack/v1"
)

func (s *service) CreateZone(ctx context.Context, in *pb.CreateZoneRequest) (*pb.CreateZoneResponse, error) {
	zone, err := s.Queries.CreateZone(ctx, db.CreateZoneParams{
		Zone:     in.Name,
		TimeZone: in.TimeZone,
	})

	if err != nil {
		return nil, QueryError{Err: err}
	}

	resp := pb.CreateZoneResponse{
		Id: zone.ID,
	}

	return &resp, nil
}

// ListZones returns a list of zones.
func (s *service) ListZones(in *pb.ListZonesRequest, out pb.StackService_ListZonesServer) error {
	zones, err := s.Queries.GetZones(context.Background(), in.ZoneGlob)
	if err != nil {
		return QueryError{Err: err}
	}

	for _, z := range zones {
		err := out.Send(&pb.ListZonesResponse{
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

func (s *service) CreateZoneAttribute(ctx context.Context, in *pb.CreateZoneAttributeRequest) (*pb.CreateZoneAttributeResponse, error) {
	attr, err := s.Queries.CreateZoneAttribute(ctx, db.CreateZoneAttributeParams{
		Zone:        in.Zone,
		Key:         in.Key,
		Value:       in.Value,
		IsProtected: bool2int(in.Protected),
	})

	if err != nil {
		return nil, QueryError{Err: err}
	}

	resp := pb.CreateZoneAttributeResponse{
		Id: attr.ID,
	}

	return &resp, nil
}
