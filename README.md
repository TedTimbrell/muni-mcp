# MUNI MCP Server

A Model Context Protocol (MCP) server for interacting with San Francisco's MUNI public transportation API.

## Overview

This project provides an MCP-compatible server that acts as a bridge between LLM applications and the SF MUNI API. It allows AI applications to query real-time public transit data, such as:

- Transit routes and route details
- Real-time arrival predictions at stops

## Getting Started

### Prerequisites

- Go 1.20 or higher

### Installation

#### Option 1: Build from Source

1. Clone the repository:
   ```
   git clone https://github.com/tedtimbrell/muni-mcp.git
   cd muni-mcp
   ```

2. Build the server for your platform:
   ```
   go build -o muni-mcp ./cmd/server
   ```

#### Option 2: Existing Cross-Platform Builds

Prebuilt binaries for multiple platforms are listed in our releases.

If you don't trust a random person on the internet you can run the multiplatform build with the following,

1. Make the build script executable:
   ```
   chmod +x scripts/build.sh
   ```

2. Run the build script:
   ```
   ./scripts/build.sh
   ```

This will generate binaries for:
- macOS (Intel and Apple Silicon)
- Linux (x86_64 and ARM64)
- Windows (x86_64)

The binaries will be available in the `build` directory with platform-specific names:
- `muni-mcp-darwin-amd64` (macOS Intel)
- `muni-mcp-darwin-arm64` (macOS Apple Silicon)
- `muni-mcp-linux-amd64` (Linux x86_64)
- `muni-mcp-linux-arm64` (Linux ARM64)
- `muni-mcp-windows-amd64.exe` (Windows x86_64)

#### MCP Config

Simply allow exectution of the bianry and place add this to your application's MCP config. Make sure to replace the command with the correct path to the binrary. `realpath <file>` is helpful for finding its absolute path.

```
{
  ...
  "mcpServers": {
    "muni-mcp": {
      "command": "/path/to/binary",
      "args": []
    }
  },
  ...
}
`

### Environment Variables

The server can be configured using the following environment variables:

- `MUNI_API_BASE_URL`: The base URL for the SF MUNI API (default: https://api.sfmta.com/v1)

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

Enable or disable caching of MUNI API responses. Defaults on to spare the poor MUNI API

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
