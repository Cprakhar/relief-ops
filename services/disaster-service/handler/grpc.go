package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
	pb "github.com/cprakhar/relief-ops/shared/proto/disaster"
	"github.com/cprakhar/relief-ops/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type GrpcHandler interface {
	ReportDisaster(ctx context.Context, req *pb.ReportDisasterRequest) (*pb.ReportDisasterResponse, error)
	ReviewDisaster(ctx context.Context, req *pb.ReviewDisasterRequest) (*pb.ReviewDisasterResponse, error)
}

// ReportDisaster handles the reporting of a new disaster.
func (h *gRPCHandler) ReportDisaster(ctx context.Context, req *pb.ReportDisasterRequest) (*pb.ReportDisasterResponse, error) {
	disaster := &types.Disaster{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Tags:        req.GetTags(),
		Location: types.Coordinates{
			Latitude:  req.GetLocation().GetLatitude(),
			Longitude: req.GetLocation().GetLongitude(),
		},
		ContributorID: req.GetContributorID(),
		ImageURLs:     req.GetImageURLs(),
	}

	// Step 1: Create the disaster in the database
	disasterID, err := h.svc.CreateDisaster(ctx, disaster)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create disaster: %v", err)
	}

	// Step 2: Try to publish to Kafka with compensation logic
	msg := &events.DisasterEventCreatedPayload{
		DisasterID:    disasterID,
		Location:      disaster.Location,
		Range:         10000,
		ContributorID: disaster.ContributorID,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal event payload: %v", err)
	}

	// Notify resource service to find resources around the disaster location
	if err := h.kafkaClient.Produce(ctx, events.ResourceCommandFind, disasterID, value); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce resource find command: %v", err)
	}
	log.Printf("Message delivered to resource-service with disasterID %s", disasterID)

	return &pb.ReportDisasterResponse{
		Id:     disasterID,
		Status: "pending",
	}, nil
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
