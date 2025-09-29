package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/services/disaster-service/types"
	"github.com/cprakhar/relief-ops/shared/events"
	"github.com/cprakhar/relief-ops/shared/messaging"
	pb "github.com/cprakhar/relief-ops/shared/proto/disaster"
	sharedTypes "github.com/cprakhar/relief-ops/shared/types"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedDisasterServiceServer
	svc         service.DisasterService
	kafkaClient *messaging.KafkaClient
}

func NewDisastergRPCHandler(srv *grpc.Server, svc service.DisasterService, kc *messaging.KafkaClient) {
	handler := &gRPCHandler{svc: svc, kafkaClient: kc}
	pb.RegisterDisasterServiceServer(srv, handler)
}

func (h *gRPCHandler) ReportDisaster(ctx context.Context, req *pb.ReportDisasterRequest) (*pb.ReportDisasterResponse, error) {
	disaster := &types.Disaster{
		ID:          uuid.New().String(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Tags:        req.GetTags(),
		Location: sharedTypes.Coordinates{
			Latitude:  req.GetLocation().GetLatitude(),
			Longitude: req.GetLocation().GetLongitude(),
		},
		ContributorID: req.GetContributorID(),
		ImageURLs:     req.GetImageURLs(),
	}

	disasterID, err := h.svc.CreateDisaster(ctx, disaster)
	if err != nil {
		return nil, err
	}

	// Step 2: Try to publish to Kafka with compensation logic
	msg := &events.DisasterEventCreatedPayload{
		DisasterID: disasterID,
		Location:   disaster.Location,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		// Compensate: Delete the created disaster
		if deleteErr := h.svc.DeleteDisaster(ctx, disasterID); deleteErr != nil {
			// Log the compensation failure but return the original error
			// In production, you'd want to send this to a dead letter queue
			// or retry mechanism
			fmt.Printf("COMPENSATION FAILED: Could not delete disaster %s: %v\n", disasterID, deleteErr)
		}
		return nil, fmt.Errorf("failed to marshal disaster event payload: %w", err)
	}

	if err := h.kafkaClient.Produce(events.ResourceCommandFind, disasterID, value); err != nil {
		// Compensate: Delete the created disaster
		if deleteErr := h.svc.DeleteDisaster(ctx, disasterID); deleteErr != nil {
			// Log the compensation failure but return the original error
			fmt.Printf("COMPENSATION FAILED: Could not delete disaster %s: %v\n", disasterID, deleteErr)
		}
		return nil, fmt.Errorf("failed to produce disaster event message: %w", err)
	}

	return &pb.ReportDisasterResponse{
		Id: disasterID,
	}, nil
}
