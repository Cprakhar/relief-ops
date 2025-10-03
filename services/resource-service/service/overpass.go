package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cprakhar/relief-ops/services/resource-service/repo"
	"github.com/cprakhar/relief-ops/shared/tools"
	"github.com/cprakhar/relief-ops/shared/types"
)

type OverpassResponse struct {
	Elements []struct {
		Type   string  `json:"type"`
		ID     int64   `json:"id"`
		Lat    float64 `json:"lat,omitempty"`
		Lon    float64 `json:"lon,omitempty"`
		Center *struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"center,omitempty"`
		Tags *struct {
			Name    string `json:"name,omitempty"`
			Amenity string `json:"amenity,omitempty"`
		} `json:"tags,omitempty"`
	} `json:"elements"`
}

type resourceService struct {
	repo repo.ResourceRepo
}

type ResourceService interface {
	SaveResources(ctx context.Context, rg int, lat, lon float64) error
	GetNearbyResources(ctx context.Context, lat, lon float64, radiusMeters int) ([]*types.Resource, error)
}

func NewResourceService(r repo.ResourceRepo) *resourceService {
	return &resourceService{repo: r}
}

func findResourcesWithinRadius(rg int, lat, lon float64) ([]*types.Resource, error) {
	overpassURL := "http://overpass-api.de/api/interpreter"

	amenities := []string{
		types.Hospital,
		types.FireStation,
		types.Police,
		types.Shelter,
		types.Pharmacy,
	}

	union := strings.Join(amenities, "|")
	query := fmt.Sprintf(`
		[out:json][timeout:30];
		(
			node["amenity"~"%s"](around:%d, %f, %f);
			way["amenity"~"%s"](around:%d, %f, %f);
			relation["amenity"~"%s"](around:%d, %f, %f);
		);
		out center;`,
		union, rg, lat, lon,
		union, rg, lat, lon,
		union, rg, lat, lon,
	)

	res, err := http.Post(overpassURL, "text/plain", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Overpass API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Overpass API response status: %d", res.StatusCode)
		return nil, fmt.Errorf("overpass API returned non-200 status: %d", res.StatusCode)
	}

	var data OverpassResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode Overpass API response: %w", err)
	}

	var resources []*types.Resource
	for _, element := range data.Elements {
		var lat, lon float64
		if element.Type == "node" {
			lat = element.Lat
			lon = element.Lon
		} else if element.Center != nil {
			lat = element.Center.Lat
			lon = element.Center.Lon
		} else {
			continue // Skip if no coordinates are available
		}

		resourceType := ""
		if element.Tags != nil {
			resourceType = element.Tags.Amenity
		}

		resource := &types.Resource{
			Name:        element.Tags.Name,
			AmenityType: resourceType,
			Location: &types.Location{
				Type:        "Point",
				Coordinates: []float64{lon, lat}, // Note: GeoJSON format is [longitude, latitude]
			},
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

func (s *resourceService) SaveResources(ctx context.Context, rg int, lat, lon float64) error {
	retryCfg := &tools.RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  time.Millisecond * 100,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2,
		Jitter:        true,
	}

	return tools.RetryWithBackoff(ctx, retryCfg, func() error {
		resources, err := findResourcesWithinRadius(rg, lat, lon)
		if err != nil {
			return err
		}
		log.Printf("Found %d resources from Overpass API", len(resources))
		log.Printf("Resources: %+v", resources)
		return s.repo.AddResources(ctx, resources)
	})
}

// GetNearbyResources retrieves resources within a certain radius (in meters) of given coordinates.
func (s *resourceService) GetNearbyResources(ctx context.Context, lat, lon float64, radiusMeters int) ([]*types.Resource, error) {
	return s.repo.GetNearbyResources(ctx, lat, lon, radiusMeters)
}
