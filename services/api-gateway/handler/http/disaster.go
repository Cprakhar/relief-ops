package http

import (
	"log"
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pbd "github.com/cprakhar/relief-ops/shared/proto/disaster"
	"github.com/cprakhar/relief-ops/shared/response"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/gin-gonic/gin"
)

type reportDisasterRequest struct {
	types.Disaster
}

// ReportDisasterHandler handles disaster reporting requests.
func ReportDisasterHandler(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req reportDisasterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
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

func GetAllDisastersHandler(ctx *gin.Context) {

}

func GetDisasterByIDHandler(ctx *gin.Context) {

}
