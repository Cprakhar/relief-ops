package http

import (
	"log"
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pbd "github.com/cprakhar/relief-ops/shared/proto/disaster"
	pbr "github.com/cprakhar/relief-ops/shared/proto/resource"
	"github.com/cprakhar/relief-ops/shared/response"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type reportDisasterRequest struct {
	types.Disaster
}

// ReportDisasterHandler handles disaster reporting requests.
func ReportDisasterHandler(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req reportDisasterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.ReportDisasterRequest{
		Title:         req.Title,
		Description:   req.Description,
		Tags:          req.Tags,
		Location:      &pbd.Coordinates{Latitude: req.Location.Latitude, Longitude: req.Location.Longitude},
		ContributorID: userID,
		ImageURLs:     req.ImageURLs,
	}

	pbRes, err := disasterClient.Client.ReportDisaster(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	responseData := struct {
		DisasterID string `json:"disaster_id"`
		Status     string `json:"status"`
	}{
		DisasterID: pbRes.GetId(),
		Status:     pbRes.GetStatus(),
	}

	ctx.JSON(http.StatusCreated, response.JSONResponse{Data: responseData})
}

func ReviewDisasterHandler(ctx *gin.Context) {
	adminID := ctx.GetString("user_id")

	disasterID := ctx.Param("id")
	decision := ctx.Query("decision") // expected values: "approve" or "reject"

	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.ReviewDisasterRequest{
		Id:      disasterID,
		AdminID: adminID,
		Status:  decision,
	}

	_, err = disasterClient.Client.ReviewDisaster(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetAllDisastersHandler retrieves all disasters.
func GetAllDisastersHandler(ctx *gin.Context) {
	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.ListDisastersRequest{
		Status: ctx.Query("status"), // Optional filter by status
	}
	pbRes, err := disasterClient.Client.ListDisasters(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	var disasters []types.Disaster
	for _, d := range pbRes.GetDisasters() {
		oid, _ := bson.ObjectIDFromHex(d.GetId())
		disaster := types.Disaster{
			ID:            oid,
			Title:         d.GetTitle(),
			Description:   d.GetDescription(),
			Tags:          d.GetTags(),
			ContributorID: d.GetContributorID(),
			CreatedAt:     d.GetCreatedAt().AsTime(),
			UpdatedAt:     d.GetUpdatedAt().AsTime(),
			ImageURLs:     d.GetImageURLs(),
			Location: types.Coordinates{
				Latitude:  d.GetLocation().GetLatitude(),
				Longitude: d.GetLocation().GetLongitude(),
			},
			Status: d.GetStatus(),
		}
		disasters = append(disasters, disaster)
	}

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: disasters})
}

// GetDisasterHandler retrieves a disaster by its ID.
func GetDisasterHandler(ctx *gin.Context) {
	disasterID := ctx.Param("id")

	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.GetDisasterRequest{Id: disasterID}
	pbRes, err := disasterClient.Client.GetDisaster(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	oid, _ := bson.ObjectIDFromHex(pbRes.GetId())
	disaster := &types.Disaster{
		ID:            oid,
		Title:         pbRes.GetTitle(),
		Description:   pbRes.GetDescription(),
		Tags:          pbRes.GetTags(),
		ContributorID: pbRes.GetContributorID(),
		CreatedAt:     pbRes.GetCreatedAt().AsTime(),
		UpdatedAt:     pbRes.GetUpdatedAt().AsTime(),
		ImageURLs:     pbRes.GetImageURLs(),
		Location: types.Coordinates{
			Latitude:  pbRes.GetLocation().GetLatitude(),
			Longitude: pbRes.GetLocation().GetLongitude(),
		},
		Status: pbRes.GetStatus(),
	}

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: disaster})
}

func GetDisasterWithResourcesHandler(ctx *gin.Context) {
	disasterID := ctx.Param("id")

	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.GetDisasterRequest{Id: disasterID}
	pbRes, err := disasterClient.Client.GetDisaster(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	oid, _ := bson.ObjectIDFromHex(pbRes.GetId())
	disaster := &types.Disaster{
		ID:            oid,
		Title:         pbRes.GetTitle(),
		Description:   pbRes.GetDescription(),
		Tags:          pbRes.GetTags(),
		ContributorID: pbRes.GetContributorID(),
		CreatedAt:     pbRes.GetCreatedAt().AsTime(),
		UpdatedAt:     pbRes.GetUpdatedAt().AsTime(),
		ImageURLs:     pbRes.GetImageURLs(),
		Location: types.Coordinates{
			Latitude:  pbRes.GetLocation().GetLatitude(),
			Longitude: pbRes.GetLocation().GetLongitude(),
		},
		Status: pbRes.GetStatus(),
	}

	resourceClient, err := grpcclient.NewResourceServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer resourceClient.Close()

	resourcesPbRes, err := resourceClient.Client.GetNearbyResources(ctx, &pbr.GetResourcesRequest{Location: &pbr.Coordinates{
		Latitude:  disaster.Location.Latitude,
		Longitude: disaster.Location.Longitude,
	}, Within: 10000}) // 10 km range
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	var resources []types.Resource
	for _, r := range resourcesPbRes.GetResources() {
		oid, _ := bson.ObjectIDFromHex(r.GetId())
		resource := types.Resource{
			ID:          oid,
			Name:        r.GetName(),
			AmenityType: r.GetAmenityType(),
			Location: &types.Location{
				Coordinates: []float64{
					r.GetLocation().Longitude,
					r.GetLocation().Latitude,
				},
			},
		}
		resources = append(resources, resource)
	}

	responseData := struct {
		Disaster  *types.Disaster  `json:"disaster"`
		Resources []types.Resource `json:"resources"`
	}{
		Disaster:  disaster,
		Resources: resources,
	}

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: responseData})
}
