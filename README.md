# MUNI MCP Server

A Model Context Protocol (MCP) server for interacting with San Francisco's MUNI public transportation API.

## Overview

This project provides an MCP-compatible server that acts as a bridge between LLM applications and the SF MUNI API. It allows AI applications to query real-time public transit data, such as:

- Transit routes
- Vehicle locations
- Arrival predictions
- Service alerts

## Getting Started

### Prerequisites

- Go 1.20 or higher
- SF MUNI API credentials (for production use)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/tedtimbrell/muni-mcp.git
   cd muni-mcp
   ```

2. Build the server:
   ```
   go build -o muni-mcp ./cmd/server
   ```

### Environment Variables

The server can be configured using the following environment variables:

- `MUNI_API_BASE_URL`: The base URL for the SF MUNI API (default: https://api.sfmta.com/v1)
- `MUNI_API_KEY`: Your API key for the SF MUNI API

### Running the Server

```
./muni-mcp
```

The server will start and listen for MCP protocol requests on stdin/stdout.

## Available Tools

### health_check

Check if the MUNI API server is healthy.

**Example:**
```json
{
  "name": "health_check"
}
```

### list_all_routes

Get a list of all MUNI routes with detailed information.

**Example:**
```json
{
  "name": "list_all_routes"
}
```

### get_route_details

Get detailed information about a specific MUNI route.

**Parameters:**
- `route_id` (string, required): ID of the route (e.g., 'N' for N-Judah)

**Example:**
```json
{
  "name": "get_route_details",
  "params": {
    "route_id": "N"
  }
}
```

### get_predictions

Get real-time arrival/departure predictions for a specific stop on a route.

**Parameters:**
- `route_id` (string, required): ID of the route (e.g., 'N' for N-Judah)
- `stop_id` (string, required): ID of the stop (e.g., '7142')

**Example:**
```json
{
  "name": "get_predictions",
  "params": {
    "route_id": "N",
    "stop_id": "7142"
  }
}
```

### toggle_cache

Enable or disable caching of MUNI API responses.

**Parameters:**
- `enabled` (boolean, required): Set to true to enable caching, false to disable

**Example:**
```json
{
  "name": "toggle_cache",
  "params": {
    "enabled": true
  }
}
```

### clear_cache

Clear the cached MUNI API responses.

**Example:**
```json
{
  "name": "clear_cache"
}
```

## Development

### Project Structure

- `cmd/server/`: Main application entry point
- `pkg/muni/`: MUNI API client implementation

### Testing

This project includes comprehensive unit tests for both the MUNI client and the MCP server handlers.

Run all tests with the following command:

```
go test ./...
```

Or use the convenience script:

```
./scripts/test.sh
```

The project also includes a GitHub Actions workflow that automatically runs tests and linting on push and pull requests.

### Adding New Tools

To add new tools to the MCP server, modify the `main.go` file in the `cmd/server` directory.

## License

MIT 