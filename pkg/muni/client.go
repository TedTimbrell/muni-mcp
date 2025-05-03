package muni

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Error constants
var (
	ErrRouteIDRequired = errors.New("route ID is required")
	ErrStopIDRequired  = errors.New("stop ID is required")
)

// cacheEntry represents a cached item with expiration
type cacheEntry struct {
	data       interface{}
	expiration time.Time
}

// isExpired checks if the cache entry has expired
func (c *cacheEntry) isExpired() bool {
	return time.Now().After(c.expiration)
}

// Cache manages cached API responses
type Cache struct {
	ttl       time.Duration
	items     map[string]cacheEntry
	mutex     sync.RWMutex
	isEnabled bool
}

// newCache creates a new cache with the given TTL
func newCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl:       ttl,
		items:     make(map[string]cacheEntry),
		isEnabled: true,
	}
}

// get retrieves an item from the cache if it exists and is not expired
func (c *Cache) get(key string, result interface{}) bool {
	if !c.isEnabled {
		return false
	}

	c.mutex.RLock()
	entry, found := c.items[key]
	c.mutex.RUnlock()

	if !found || entry.isExpired() {
		return false
	}

	// Copy the cached data to the result
	data, err := json.Marshal(entry.data)
	if err != nil {
		return false
	}

	if err := json.Unmarshal(data, result); err != nil {
		return false
	}

	return true
}

// set adds or updates an item in the cache
func (c *Cache) set(key string, data interface{}) {
	if !c.isEnabled {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = cacheEntry{
		data:       data,
		expiration: time.Now().Add(c.ttl),
	}
}

// clear removes all items from the cache
func (c *Cache) clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]cacheEntry)
}

// enable turns on caching
func (c *Cache) enable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isEnabled = true
}

// disable turns off caching
func (c *Cache) disable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isEnabled = false
}

// Client represents a client for the SF MUNI API
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	cache      *Cache
}

// ClientOption is a functional option for configuring the client
type ClientOption func(*Client)

// WithCacheTTL sets the cache time-to-live duration
func WithCacheTTL(ttl time.Duration) ClientOption {
	return func(c *Client) {
		c.cache = newCache(ttl)
	}
}

// WithoutCache disables caching
func WithoutCache() ClientOption {
	return func(c *Client) {
		c.cache.disable()
	}
}

// NewClient creates a new MUNI API client
func NewClient(baseURL, apiKey string, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
		cache:   newCache(5 * time.Minute), // Default cache TTL
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ClearCache clears all cached responses
func (c *Client) ClearCache() {
	c.cache.clear()
}

// EnableCache enables caching
func (c *Client) EnableCache() {
	c.cache.enable()
}

// DisableCache disables caching
func (c *Client) DisableCache() {
	c.cache.disable()
}

// GetAllRoutes fetches all available MUNI routes with detailed information
func (c *Client) GetAllRoutes(ctx context.Context) ([]RouteInfo, error) {
	cacheKey := "all_routes"

	// Try to get from cache first
	var routes []RouteInfo
	if c.cache.get(cacheKey, &routes) {
		return routes, nil
	}

	url := fmt.Sprintf("%s/v2.0/riders/agencies/sfmta-cis/routes", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&routes); err != nil {
		return nil, err
	}

	// Cache the response
	c.cache.set(cacheKey, routes)

	return routes, nil
}

// GetRouteDetails fetches detailed information for a specific route
func (c *Client) GetRouteDetails(ctx context.Context, routeID string) (*RouteDetails, error) {
	if routeID == "" {
		return nil, ErrRouteIDRequired
	}

	cacheKey := fmt.Sprintf("route_details:%s", routeID)

	// Try to get from cache first
	var routeDetails RouteDetails
	if c.cache.get(cacheKey, &routeDetails) {
		return &routeDetails, nil
	}

	url := fmt.Sprintf("%s/v2.0/riders/agencies/sfmta-cis/routes/%s", c.baseURL, routeID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&routeDetails); err != nil {
		return nil, err
	}

	// Cache the response
	c.cache.set(cacheKey, routeDetails)

	return &routeDetails, nil
}

// GetPredictions fetches real-time predictions for a specific stop on a route
func (c *Client) GetPredictions(ctx context.Context, routeID, stopID string) ([]Prediction, error) {
	if routeID == "" {
		return nil, ErrRouteIDRequired
	}

	if stopID == "" {
		return nil, ErrStopIDRequired
	}

	url := fmt.Sprintf("%s/v2.0/riders/agencies/sfmta-cis/nstops/%s:%s/predictions", c.baseURL, routeID, stopID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var predictionResponse []PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predictionResponse); err != nil {
		return nil, err
	}

	// If there are no prediction responses or no values in the first response, return empty predictions
	if len(predictionResponse) == 0 || len(predictionResponse[0].Values) == 0 {
		return []Prediction{}, nil
	}

	// Convert prediction response to predictions
	predictions := make([]Prediction, len(predictionResponse[0].Values))
	for i, val := range predictionResponse[0].Values {
		predictions[i] = Prediction{
			VehicleID:       val.VehicleID,
			Minutes:         val.Minutes,
			Direction:       val.Direction.Name,
			DestinationName: val.Direction.DestinationName,
			Timestamp:       time.Unix(val.Timestamp/1000, 0),
			VehicleType:     val.VehicleType,
			IsDeparture:     val.IsDeparture,
		}
	}

	return predictions, nil
}

// RouteInfo represents basic information about a MUNI route from the API
type RouteInfo struct {
	ID          string `json:"id"`
	Rev         int    `json:"rev"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	TextColor   string `json:"textColor"`
	Hidden      bool   `json:"hidden"`
	Timestamp   string `json:"timestamp"`
}

// BoundingBox represents the geographical bounds of a route
type BoundingBox struct {
	LatMin float64 `json:"latMin"`
	LatMax float64 `json:"latMax"`
	LonMin float64 `json:"lonMin"`
	LonMax float64 `json:"lonMax"`
}

// Stop represents a transit stop
type Stop struct {
	ID                      string   `json:"id"`
	Lat                     float64  `json:"lat"`
	Lon                     float64  `json:"lon"`
	Name                    string   `json:"name"`
	Code                    string   `json:"code,omitempty"`
	Hidden                  bool     `json:"hidden"`
	ShowDestinationSelector bool     `json:"showDestinationSelector"`
	Directions              []string `json:"directions"`
}

// Direction represents a direction of travel on a route
type Direction struct {
	ID        string   `json:"id"`
	ShortName string   `json:"shortName"`
	Name      string   `json:"name"`
	UseForUI  bool     `json:"useForUi"`
	Stops     []string `json:"stops"`
}

// Path represents a path segment on a route
type Path struct {
	ID     string      `json:"id"`
	Points []PathPoint `json:"points"`
}

// PathPoint represents a single point on a path
type PathPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// RouteDetails represents detailed information about a MUNI route
type RouteDetails struct {
	ID          string      `json:"id"`
	Rev         int         `json:"rev"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Color       string      `json:"color"`
	TextColor   string      `json:"textColor"`
	Hidden      bool        `json:"hidden"`
	BoundingBox BoundingBox `json:"boundingBox"`
	Stops       []Stop      `json:"stops"`
	Directions  []Direction `json:"directions"`
	Paths       []Path      `json:"paths"`
	Timestamp   string      `json:"timestamp"`
}

// PredictionDirection represents information about the direction of a prediction
type PredictionDirection struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DestinationName string `json:"destinationName"`
}

// PredictionValue represents a single prediction value
type PredictionValue struct {
	Timestamp             int64               `json:"timestamp"`
	Minutes               int                 `json:"minutes"`
	AffectedByLayover     bool                `json:"affectedByLayover"`
	IsDeparture           bool                `json:"isDeparture"`
	OccupancyStatus       int                 `json:"occupancyStatus"`
	OccupancyDescription  string              `json:"occupancyDescription"`
	VehiclesInConsist     int                 `json:"vehiclesInConsist"`
	LinkedVehicleIds      string              `json:"linkedVehicleIds"`
	VehicleID             string              `json:"vehicleId"`
	VehicleType           string              `json:"vehicleType"`
	Direction             PredictionDirection `json:"direction"`
	TripID                string              `json:"tripId"`
	Delay                 int                 `json:"delay"`
	PredUsingNavigationTm bool                `json:"predUsingNavigationTm"`
	Departure             bool                `json:"departure"`
}

// PredictionRoute represents the route of a prediction
type PredictionRoute struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	TextColor   string `json:"textColor"`
	Hidden      bool   `json:"hidden"`
}

// PredictionStop represents the stop of a prediction
type PredictionStop struct {
	ID                      string  `json:"id"`
	Lat                     float64 `json:"lat"`
	Lon                     float64 `json:"lon"`
	Name                    string  `json:"name"`
	Code                    string  `json:"code,omitempty"`
	Hidden                  bool    `json:"hidden"`
	ShowDestinationSelector bool    `json:"showDestinationSelector"`
	Route                   string  `json:"route"`
}

// Agency represents the transit agency information
type Agency struct {
	Rev       int    `json:"rev"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	// Additional fields omitted for brevity
}

// PredictionResponse represents the response from the predictions API
type PredictionResponse struct {
	ServerTimestamp int64             `json:"serverTimestamp"`
	NxbsRedirectURL string            `json:"nxbs2RedirectUrl"`
	Agency          Agency            `json:"agency"`
	Route           PredictionRoute   `json:"route"`
	Stop            PredictionStop    `json:"stop"`
	Values          []PredictionValue `json:"values"`
}

// Prediction represents a simplified prediction for a vehicle arrival/departure
type Prediction struct {
	VehicleID       string    `json:"vehicle_id"`
	Minutes         int       `json:"minutes"`
	Direction       string    `json:"direction"`
	DestinationName string    `json:"destination_name"`
	Timestamp       time.Time `json:"timestamp"`
	VehicleType     string    `json:"vehicle_type"`
	IsDeparture     bool      `json:"is_departure"`
}
