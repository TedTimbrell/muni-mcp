package muni

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://test-api.example.com"
	apiKey := "test-api-key"

	client := NewClient(baseURL, apiKey)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL to be %s, got %s", baseURL, client.baseURL)
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %s, got %s", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be %s, got %s", 10*time.Second, client.httpClient.Timeout)
	}
}

func TestGetRoutes(t *testing.T) {
	client := NewClient("https://test-api.example.com", "test-api-key")

	routes, err := client.GetRoutes(context.Background())

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

	if routes[0].Name != "N-Judah" {
		t.Errorf("Expected route name to be N-Judah, got %s", routes[0].Name)
	}

	if routes[0].Type != "light_rail" {
		t.Errorf("Expected route type to be light_rail, got %s", routes[0].Type)
	}

	// Check second route
	if routes[1].ID != "J" {
		t.Errorf("Expected route ID to be J, got %s", routes[1].ID)
	}

	if routes[1].Name != "J-Church" {
		t.Errorf("Expected route name to be J-Church, got %s", routes[1].Name)
	}

	if routes[1].Type != "light_rail" {
		t.Errorf("Expected route type to be light_rail, got %s", routes[1].Type)
	}
}

func TestGetVehicleLocations(t *testing.T) {
	client := NewClient("https://test-api.example.com", "test-api-key")

	// Test with valid route ID
	routeID := "N"
	vehicles, err := client.GetVehicleLocations(context.Background(), routeID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(vehicles) != 1 {
		t.Errorf("Expected 1 vehicle, got %d", len(vehicles))
	}

	vehicle := vehicles[0]

	if vehicle.VehicleID != "1234" {
		t.Errorf("Expected vehicle ID to be 1234, got %s", vehicle.VehicleID)
	}

	if vehicle.RouteID != routeID {
		t.Errorf("Expected route ID to be %s, got %s", routeID, vehicle.RouteID)
	}

	if vehicle.Latitude != 37.7749 {
		t.Errorf("Expected latitude to be 37.7749, got %f", vehicle.Latitude)
	}

	if vehicle.Longitude != -122.4194 {
		t.Errorf("Expected longitude to be -122.4194, got %f", vehicle.Longitude)
	}

	if vehicle.Heading != 90 {
		t.Errorf("Expected heading to be 90, got %d", vehicle.Heading)
	}

	if vehicle.Speed != 15.5 {
		t.Errorf("Expected speed to be 15.5, got %f", vehicle.Speed)
	}

	// Test with empty route ID
	_, err = client.GetVehicleLocations(context.Background(), "")

	if err == nil {
		t.Error("Expected error for empty route ID, got nil")
	}
}
