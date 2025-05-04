package muni

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://test.com"
	client := NewClient(baseURL)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL to be %s, got %s", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if client.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestClientOptions(t *testing.T) {
	// Test WithCacheTTL option
	ttl := 10 * time.Minute
	client := NewClient("http://test.com", WithCacheTTL(ttl))
	if client.cache.ttl != ttl {
		t.Errorf("Expected cache TTL to be %v, got %v", ttl, client.cache.ttl)
	}

	// Test WithoutCache option
	client = NewClient("http://test.com", WithoutCache())
	if client.cache.isEnabled {
		t.Error("Expected cache to be disabled")
	}
}

func TestGetAllRoutes(t *testing.T) {
	server := mockServer(mockRoutesResponse)
	defer server.Close()

	client := NewClient(server.URL)
	routes, err := client.GetAllRoutes(context.Background())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	// Check first route
	if routes[0].ID != "N" {
		t.Errorf("Expected route ID to be N, got %s", routes[0].ID)
	}

	if routes[0].Title != "N-Judah" {
		t.Errorf("Expected route title to be N-Judah, got %s", routes[0].Title)
	}

	// Check second route
	if routes[1].ID != "J" {
		t.Errorf("Expected route ID to be J, got %s", routes[1].ID)
	}

	if routes[1].Title != "J-Church" {
		t.Errorf("Expected route title to be J-Church, got %s", routes[1].Title)
	}
}

func TestGetRouteDetails(t *testing.T) {
	server := mockServer(mockRouteDetailsResponse)
	defer server.Close()

	client := NewClient(server.URL)

	// Test with valid route ID
	routeID := "N"
	details, err := client.GetRouteDetails(context.Background(), routeID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if details.ID != routeID {
		t.Errorf("Expected route ID to be %s, got %s", routeID, details.ID)
	}

	if len(details.Stops) != 1 {
		t.Errorf("Expected 1 stop, got %d", len(details.Stops))
	}

	if details.Stops[0].ID != "1234" {
		t.Errorf("Expected stop ID to be 1234, got %s", details.Stops[0].ID)
	}

	// Test with empty route ID
	_, err = client.GetRouteDetails(context.Background(), "")
	if err != ErrRouteIDRequired {
		t.Errorf("Expected ErrRouteIDRequired, got %v", err)
	}
}

func TestGetPredictions(t *testing.T) {
	server := mockServer(mockPredictionsResponse)
	defer server.Close()

	client := NewClient(server.URL)

	// Test with valid route and stop IDs
	routeID := "N"
	stopID := "1234"
	predictions, err := client.GetPredictions(context.Background(), routeID, stopID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(predictions) == 0 {
		t.Error("Expected predictions, got none")
	}

	prediction := predictions[0]
	if prediction.VehicleID != "1234" {
		t.Errorf("Expected vehicle ID to be 1234, got %s", prediction.VehicleID)
	}

	if prediction.Minutes != 5 {
		t.Errorf("Expected minutes to be 5, got %d", prediction.Minutes)
	}

	if prediction.Direction != "Inbound" {
		t.Errorf("Expected direction to be Inbound, got %s", prediction.Direction)
	}

	// Test with empty route ID
	_, err = client.GetPredictions(context.Background(), "", stopID)
	if err != ErrRouteIDRequired {
		t.Errorf("Expected ErrRouteIDRequired, got %v", err)
	}

	// Test with empty stop ID
	_, err = client.GetPredictions(context.Background(), routeID, "")
	if err != ErrStopIDRequired {
		t.Errorf("Expected ErrStopIDRequired, got %v", err)
	}
}

func TestCacheOperations(t *testing.T) {
	client := NewClient("https://test-api.example.com")

	// Test cache enable/disable
	client.DisableCache()
	if client.cache.isEnabled {
		t.Error("Expected cache to be disabled")
	}

	client.EnableCache()
	if !client.cache.isEnabled {
		t.Error("Expected cache to be enabled")
	}

	// Test cache clear
	client.cache.set("test", "data")
	client.ClearCache()
	var result string
	if client.cache.get("test", &result) {
		t.Error("Expected cache to be empty after clear")
	}
}
