package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tedtimbrell/muni-mcp/pkg/muni"
)

func TestHealthCheckHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	// Test success case
	mockClient.GetAllRoutesFunc = func(ctx context.Context) ([]muni.RouteInfo, error) {
		return []muni.RouteInfo{
			{ID: "N", Title: "N Judah"},
		}, nil
	}

	result, err := healthCheckHandler(context.Background(), mcp.CallToolRequest{})

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	expectedText := "SF MUNI API server is healthy and running!"
	if textContent.Text != expectedText {
		t.Errorf("Expected text to be '%s', got '%s'", expectedText, textContent.Text)
	}

	// Test failure case
	mockClient.GetAllRoutesFunc = func(ctx context.Context) ([]muni.RouteInfo, error) {
		return nil, errors.New("API connection error")
	}

	result, err = healthCheckHandler(context.Background(), mcp.CallToolRequest{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true for health check failure")
	}
}

func TestListAllRoutesHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	// Test success case
	mockClient.GetAllRoutesFunc = func(ctx context.Context) ([]muni.RouteInfo, error) {
		return []muni.RouteInfo{
			{
				ID:          "N",
				Rev:         1069,
				Title:       "N Judah",
				Description: "Weekdays 6am-12 midnight Weekends 8am-12 midnight",
				Color:       "005b95",
				TextColor:   "ffffff",
				Hidden:      false,
				Timestamp:   "2025-04-26T10:31:08Z",
			},
			{
				ID:          "J",
				Rev:         1069,
				Title:       "J Church",
				Description: "5am-12 midnight daily",
				Color:       "a96614",
				TextColor:   "ffffff",
				Hidden:      false,
				Timestamp:   "2025-04-26T10:31:08Z",
			},
		}, nil
	}

	result, err := listAllRoutesHandler(context.Background(), mcp.CallToolRequest{})

	// Assert success case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse the JSON to verify content
	var routes []muni.RouteInfo
	if err := json.Unmarshal([]byte(textContent.Text), &routes); err != nil {
		t.Fatalf("Failed to unmarshal routes: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	if routes[0].ID != "N" || routes[1].ID != "J" {
		t.Errorf("Route IDs don't match expected values")
	}

	// Test error case
	mockClient.GetAllRoutesFunc = func(ctx context.Context) ([]muni.RouteInfo, error) {
		return nil, errors.New("API error")
	}

	result, err = listAllRoutesHandler(context.Background(), mcp.CallToolRequest{})

	// Assert error case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}
}

func TestGetRouteDetailsHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	// Test success case
	mockClient.GetRouteDetailsFunc = func(ctx context.Context, routeID string) (*muni.RouteDetails, error) {
		if routeID == "" {
			return nil, muni.ErrRouteIDRequired
		}
		return &muni.RouteDetails{
			ID:          routeID,
			Rev:         1069,
			Title:       routeID + " Test Route",
			Description: "Test route description",
			Color:       "005b95",
			TextColor:   "ffffff",
		}, nil
	}

	routeID := "N"
	request := mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"route_id": routeID,
	}

	result, err := getRouteDetailsHandler(context.Background(), request)

	// Assert success case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse the JSON to verify content
	var details muni.RouteDetails
	if err := json.Unmarshal([]byte(textContent.Text), &details); err != nil {
		t.Fatalf("Failed to unmarshal route details: %v", err)
	}

	if details.ID != routeID {
		t.Errorf("Expected route ID to be %s, got %s", routeID, details.ID)
	}

	// Test missing route_id parameter
	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{}

	result, err = getRouteDetailsHandler(context.Background(), request)

	// Assert missing route_id parameter case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	// Test API error case
	mockClient.GetRouteDetailsFunc = func(ctx context.Context, routeID string) (*muni.RouteDetails, error) {
		return nil, errors.New("API error")
	}

	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"route_id": routeID,
	}

	result, err = getRouteDetailsHandler(context.Background(), request)

	// Assert API error case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}
}

func TestGetPredictionsHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	// Test success case
	mockClient.GetPredictionsFunc = func(ctx context.Context, routeID, stopID string) ([]muni.Prediction, error) {
		if routeID == "" {
			return nil, muni.ErrRouteIDRequired
		}
		if stopID == "" {
			return nil, muni.ErrStopIDRequired
		}
		return []muni.Prediction{
			{
				VehicleID:       "51",
				Minutes:         9,
				Direction:       "Market & California",
				DestinationName: "Market & California",
				Timestamp:       time.Now().Add(9 * time.Minute),
				VehicleType:     "Cable Car_CABLECAR",
				IsDeparture:     true,
			},
		}, nil
	}

	routeID := "CA"
	stopID := "7142"
	request := mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"route_id": routeID,
		"stop_id":  stopID,
	}

	result, err := getPredictionsHandler(context.Background(), request)

	// Assert success case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	// Parse the JSON to verify content
	var predictions []muni.Prediction
	if err := json.Unmarshal([]byte(textContent.Text), &predictions); err != nil {
		t.Fatalf("Failed to unmarshal predictions: %v", err)
	}

	if len(predictions) != 1 {
		t.Errorf("Expected 1 prediction, got %d", len(predictions))
	}

	// Test missing route_id parameter
	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"stop_id": stopID,
	}

	result, err = getPredictionsHandler(context.Background(), request)

	// Assert missing route_id parameter case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	// Test missing stop_id parameter
	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"route_id": routeID,
	}

	result, err = getPredictionsHandler(context.Background(), request)

	// Assert missing stop_id parameter case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	// Test API error case
	mockClient.GetPredictionsFunc = func(ctx context.Context, routeID, stopID string) ([]muni.Prediction, error) {
		return nil, errors.New("API error")
	}

	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"route_id": routeID,
		"stop_id":  stopID,
	}

	result, err = getPredictionsHandler(context.Background(), request)

	// Assert API error case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}
}

func TestClearCacheHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	clearCalled := false
	mockClient.ClearCacheFunc = func() {
		clearCalled = true
	}

	// Test
	result, err := clearCacheHandler(context.Background(), mcp.CallToolRequest{})

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !clearCalled {
		t.Error("Expected ClearCache to be called")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	expectedText := "MUNI API cache has been cleared"
	if textContent.Text != expectedText {
		t.Errorf("Expected text to be '%s', got '%s'", expectedText, textContent.Text)
	}
}

func TestToggleCacheHandler(t *testing.T) {
	// Setup
	originalClient := muniClient
	defer func() { muniClient = originalClient }()

	mockClient := muni.NewMockClient()
	muniClient = mockClient

	enableCalled := false
	disableCalled := false

	mockClient.EnableCacheFunc = func() {
		enableCalled = true
	}

	mockClient.DisableCacheFunc = func() {
		disableCalled = true
	}

	// Test enable case
	request := mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"enabled": true,
	}

	result, err := toggleCacheHandler(context.Background(), request)

	// Assert enable case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !enableCalled {
		t.Error("Expected EnableCache to be called")
	}

	if disableCalled {
		t.Error("DisableCache should not have been called")
	}

	// Test disable case
	enableCalled = false
	disableCalled = false

	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{
		"enabled": false,
	}

	result, err = toggleCacheHandler(context.Background(), request)

	// Assert disable case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if enableCalled {
		t.Error("EnableCache should not have been called")
	}

	if !disableCalled {
		t.Error("Expected DisableCache to be called")
	}

	// Test missing parameter case
	request = mcp.CallToolRequest{}
	request.Params.Arguments = map[string]interface{}{}

	result, err = toggleCacheHandler(context.Background(), request)

	// Assert missing parameter case
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}
}
