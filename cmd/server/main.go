package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedtimbrell/muni-mcp/pkg/muni"
)

// MuniClient is the interface for interacting with the MUNI API
type MuniClient interface {
	GetAllRoutes(ctx context.Context) ([]muni.RouteInfo, error)
	GetRouteDetails(ctx context.Context, routeID string) (*muni.RouteDetails, error)
	GetPredictions(ctx context.Context, routeID, stopID string) ([]muni.Prediction, error)
	ClearCache()
	EnableCache()
	DisableCache()
}

var muniClient MuniClient

func main() {
	// Initialize MUNI client
	// In a production environment, these would come from environment variables
	baseURL := os.Getenv("MUNI_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.prd-1.iq.live.umoiq.com" // Default URL based on HAR files
	}

	apiKey := os.Getenv("MUNI_API_KEY")
	cacheTTL := 5 * time.Minute

	// Parse cache TTL from environment if present
	if ttlStr := os.Getenv("MUNI_CACHE_TTL"); ttlStr != "" {
		if parsedTTL, err := time.ParseDuration(ttlStr); err == nil {
			cacheTTL = parsedTTL
		}
	}

	muniClient = muni.NewClient(baseURL, apiKey, muni.WithCacheTTL(cacheTTL))

	// Create MCP server
	s := server.NewMCPServer(
		"SF MUNI API Server",
		"0.2.0",
	)

	// Add a simple health check tool
	healthTool := mcp.NewTool("health_check",
		mcp.WithDescription("Check if the MUNI API server is healthy"),
	)

	// Add route listing tool
	allRoutesTool := mcp.NewTool("list_all_routes",
		mcp.WithDescription("Get a list of all MUNI routes with detailed information"),
	)

	// Add route details tool
	routeDetailsTool := mcp.NewTool("get_route_details",
		mcp.WithDescription("Get detailed information about a specific MUNI route"),
		mcp.WithString("route_id",
			mcp.Required(),
			mcp.Description("ID of the route (e.g., 'N' for N-Judah)"),
		),
	)

	// Add predictions tool
	predictionsTool := mcp.NewTool("get_predictions",
		mcp.WithDescription("Get real-time arrival/departure predictions for a specific stop on a route"),
		mcp.WithString("route_id",
			mcp.Required(),
			mcp.Description("ID of the route (e.g., 'N' for N-Judah)"),
		),
		mcp.WithString("stop_id",
			mcp.Required(),
			mcp.Description("ID of the stop (e.g., '7142')"),
		),
	)

	// Add cache management tools
	clearCacheTool := mcp.NewTool("clear_cache",
		mcp.WithDescription("Clear the cached MUNI API responses"),
	)

	toggleCacheTool := mcp.NewTool("toggle_cache",
		mcp.WithDescription("Enable or disable caching of MUNI API responses"),
		mcp.WithBoolean("enabled",
			mcp.Required(),
			mcp.Description("Set to true to enable caching, false to disable"),
		),
	)

	// Add tool handlers
	s.AddTool(healthTool, healthCheckHandler)
	s.AddTool(allRoutesTool, listAllRoutesHandler)
	s.AddTool(routeDetailsTool, getRouteDetailsHandler)
	s.AddTool(predictionsTool, getPredictionsHandler)
	s.AddTool(clearCacheTool, clearCacheHandler)
	s.AddTool(toggleCacheTool, toggleCacheHandler)

	// Start the stdio server
	log.Println("Starting SF MUNI MCP server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// Create a helper function for JSON content
func newJSONToolResult(data interface{}) (*mcp.CallToolResult, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal to JSON: %v", err)), nil
	}

	// Since we can't set MIME type in the current version of MCP Go, we'll just use text content
	textContent := mcp.NewTextContent(string(jsonData))
	textContent.Type = "text"

	return &mcp.CallToolResult{
		Content: []mcp.Content{textContent},
	}, nil
}

func healthCheckHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Check if we can get a list of routes as a basic connectivity test
	_, err := muniClient.GetAllRoutes(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("MUNI API health check failed: %v", err)), nil
	}

	return mcp.NewToolResultText("SF MUNI API server is healthy and running!"), nil
}

func listAllRoutesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	routes, err := muniClient.GetAllRoutes(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch routes: %v", err)), nil
	}

	return newJSONToolResult(routes)
}

func getRouteDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	routeID, ok := request.Params.Arguments["route_id"].(string)
	if !ok {
		return mcp.NewToolResultError("route_id must be a string"), nil
	}

	details, err := muniClient.GetRouteDetails(ctx, routeID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch route details: %v", err)), nil
	}

	return newJSONToolResult(details)
}

func getPredictionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	routeID, ok := request.Params.Arguments["route_id"].(string)
	if !ok {
		return mcp.NewToolResultError("route_id must be a string"), nil
	}

	stopID, ok := request.Params.Arguments["stop_id"].(string)
	if !ok {
		return mcp.NewToolResultError("stop_id must be a string"), nil
	}

	predictions, err := muniClient.GetPredictions(ctx, routeID, stopID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch predictions: %v", err)), nil
	}

	return newJSONToolResult(predictions)
}

func clearCacheHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	muniClient.ClearCache()
	return mcp.NewToolResultText("MUNI API cache has been cleared"), nil
}

func toggleCacheHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	enabled, ok := request.Params.Arguments["enabled"].(bool)
	if !ok {
		return mcp.NewToolResultError("enabled must be a boolean"), nil
	}

	if enabled {
		muniClient.EnableCache()
		return mcp.NewToolResultText("MUNI API caching is now enabled"), nil
	} else {
		muniClient.DisableCache()
		return mcp.NewToolResultText("MUNI API caching is now disabled"), nil
	}
}
