package handler

import (
	"context"

	"github.com/cprakhar/relief-ops/services/resource-service/service"
	pb "github.com/cprakhar/relief-ops/shared/proto/resource"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedResourceServiceServer
	svc service.ResourceService
}

type GrpcHandler interface {
	GetNearbyResources(ctx context.Context, req *pb.GetResourcesRequest) (*pb.GetResourcesResponse, error)
}

// NewResourcegRPCHandler registers the gRPC handler for the ResourceService.
func NewResourcegRPCHandler(srv *grpc.Server, svc service.ResourceService) {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterResourceServiceServer(srv, handler)
}

// GetNearbyResources handles requests to fetch nearby resources based on given coordinates and radius.
func (h *gRPCHandler) GetNearbyResources(ctx context.Context, req *pb.GetResourcesRequest) (*pb.GetResourcesResponse, error) {
	within := req.GetWithin()
	lon, lat := req.GetLocation().GetLongitude(), req.GetLocation().GetLatitude()

	resources, err := h.svc.GetNearbyResources(ctx, lat, lon, int(within))
	if err != nil {
		return nil, err
	}

	var pbResources []*pb.Resource
	for _, r := range resources {
		pbResource := &pb.Resource{
			Name:        r.Name,
			AmenityType: r.AmenityType,
			Location: &pb.Coordinates{
				Latitude:  r.Location.Coordinates[1], // GeoJSON format is [longitude, latitude]
				Longitude: r.Location.Coordinates[0],
			},
		}
		pbResources = append(pbResources, pbResource)
	}

	return &pb.GetResourcesResponse{
		Resources: pbResources,
	}, nil
}
