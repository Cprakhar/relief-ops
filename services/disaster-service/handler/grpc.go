package handler

import (
	"context"
	"encoding/json"

	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	pb "github.com/cprakhar/relief-ops/shared/proto/disaster"
	"github.com/cprakhar/relief-ops/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCHandler struct {
	pb.UnimplementedDisasterServiceServer
	svc         service.DisasterService
	kafkaClient *messaging.KafkaClient
}

// NewDisastergRPCHandler registers the gRPC handler for the DisasterService.
func NewDisastergRPCHandler(srv *grpc.Server, svc service.DisasterService, kc *messaging.KafkaClient) {
	handler := &gRPCHandler{svc: svc, kafkaClient: kc}
	pb.RegisterDisasterServiceServer(srv, handler)
}

// GrpcHandler defines the interface for gRPC handler methods.
type GrpcHandler interface {
	ListDisasters(ctx context.Context, req *pb.ListDisastersRequest) (*pb.ListDisastersResponse, error)
	GetDisaster(ctx context.Context, req *pb.GetDisasterRequest) (*pb.GetDisasterResponse, error)
	ReportDisaster(ctx context.Context, req *pb.ReportDisasterRequest) (*pb.ReportDisasterResponse, error)
	ReviewDisaster(ctx context.Context, req *pb.ReviewDisasterRequest) (*pb.ReviewDisasterResponse, error)
}

// ReportDisaster handles the reporting of a new disaster.
func (h *gRPCHandler) ReportDisaster(ctx context.Context, req *pb.ReportDisasterRequest) (*pb.ReportDisasterResponse, error) {
	logger := logs.L()

	disaster := &types.Disaster{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Tags:        req.GetTags(),
		Location: types.Coordinates{
			Latitude:  req.GetLocation().GetLatitude(),
			Longitude: req.GetLocation().GetLongitude(),
		},
		VolunteerID: req.GetVolunteerID(),
		ImageURLs:   req.GetImageURLs(),
	}

	// Step 1: Create the disaster in the database
	disasterID, err := h.svc.CreateDisaster(ctx, disaster)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create disaster: %v", err)
	}

	// Step 2: Try to publish to Kafka with compensation logic
	msg := &events.DisasterEventCreatedPayload{
		DisasterID:  disasterID,
		Location:    disaster.Location,
		Range:       10000,
		VolunteerID: disaster.VolunteerID,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal event payload: %v", err)
	}

	// Notify resource service to find resources around the disaster location
	logger.Infow("Notifying resource service to find resources", "disaster_id", disasterID, "location", disaster.Location, "range", 10000)
	if err := h.kafkaClient.Produce(ctx, events.ResourceCommandFind, disasterID, value); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce resource find command: %v", err)
	}

	return &pb.ReportDisasterResponse{
		Id:     disasterID,
		Status: "pending",
	}, nil
}

// GetDisaster retrieves a disaster by its ID.
func (h *gRPCHandler) GetDisaster(ctx context.Context, req *pb.GetDisasterRequest) (*pb.GetDisasterResponse, error) {
	disasterID := req.GetId()
	disaster, err := h.svc.GetDisaster(ctx, disasterID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get disaster: %v", err)
	}

	return &pb.GetDisasterResponse{
		Id:          disaster.ID.Hex(),
		Title:       disaster.Title,
		Description: disaster.Description,
		Tags:        disaster.Tags,
		VolunteerID: disaster.VolunteerID,
		CreatedAt:   timestamppb.New(disaster.CreatedAt),
		UpdatedAt:   timestamppb.New(disaster.UpdatedAt),
		ImageURLs:   disaster.ImageURLs,
		Location: &pb.Coordinates{
			Latitude:  disaster.Location.Latitude,
			Longitude: disaster.Location.Longitude,
		},
		Status: disaster.Status,
	}, nil
}

// ListDisasters retrieves all disasters, optionally filtered by status.
func (h *gRPCHandler) ListDisasters(ctx context.Context, req *pb.ListDisastersRequest) (*pb.ListDisastersResponse, error) {
	disasters, err := h.svc.GetAllDisasters(ctx, req.GetStatus())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list disasters: %v", err)
	}

	var pbDisasters []*pb.GetDisasterResponse
	for _, d := range disasters {
		pbDisasters = append(pbDisasters, &pb.GetDisasterResponse{
			Id:          d.ID.Hex(),
			Title:       d.Title,
			Description: d.Description,
			Tags:        d.Tags,
			VolunteerID: d.VolunteerID,
			CreatedAt:   timestamppb.New(d.CreatedAt),
			UpdatedAt:   timestamppb.New(d.UpdatedAt),
			ImageURLs:   d.ImageURLs,
			Location: &pb.Coordinates{
				Latitude:  d.Location.Latitude,
				Longitude: d.Location.Longitude,
			},
			Status: d.Status,
		})
	}

	return &pb.ListDisastersResponse{Disasters: pbDisasters}, nil
}

// ReviewDisaster handles the review of a reported disaster.
func (h *gRPCHandler) ReviewDisaster(ctx context.Context, req *pb.ReviewDisasterRequest) (*pb.ReviewDisasterResponse, error) {
	if err := h.svc.UpdateStatus(ctx, req.GetId(), req.GetStatus()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update disaster status: %v", err)
	}

	return &pb.ReviewDisasterResponse{
		Id:     req.GetId(),
		Status: req.GetStatus(),
	}, nil
}
