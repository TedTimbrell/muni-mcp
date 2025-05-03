package muni

import (
	"context"
	"time"
)

// MockClient is a mock implementation of the MUNI client for testing
type MockClient struct {
	GetAllRoutesFunc    func(ctx context.Context) ([]RouteInfo, error)
	GetRouteDetailsFunc func(ctx context.Context, routeID string) (*RouteDetails, error)
	GetPredictionsFunc  func(ctx context.Context, routeID, stopID string) ([]Prediction, error)
	ClearCacheFunc      func()
	EnableCacheFunc     func()
	DisableCacheFunc    func()
}

// Ensure MockClient implements required interface
var _ interface {
	GetAllRoutes(ctx context.Context) ([]RouteInfo, error)
	GetRouteDetails(ctx context.Context, routeID string) (*RouteDetails, error)
	GetPredictions(ctx context.Context, routeID, stopID string) ([]Prediction, error)
	ClearCache()
	EnableCache()
	DisableCache()
} = (*MockClient)(nil)

// NewMockClient creates a new mock MUNI client with default implementations
func NewMockClient() *MockClient {
	return &MockClient{
		GetAllRoutesFunc: func(ctx context.Context) ([]RouteInfo, error) {
			return []RouteInfo{
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
		},
		GetRouteDetailsFunc: func(ctx context.Context, routeID string) (*RouteDetails, error) {
			if routeID == "" {
				return nil, ErrRouteIDRequired
			}

			return &RouteDetails{
				ID:          routeID,
				Rev:         1069,
				Title:       routeID + " Test Route",
				Description: "Test route description",
				Color:       "005b95",
				TextColor:   "ffffff",
				Hidden:      false,
				BoundingBox: BoundingBox{
					LatMin: 37.7904099,
					LatMax: 37.7936799,
					LonMin: -122.42224,
					LonMax: -122.39637,
				},
				Stops: []Stop{
					{
						ID:                      "3860",
						Lat:                     37.7936799,
						Lon:                     -122.39637,
						Name:                    "Test Stop 1",
						Code:                    "13860",
						Hidden:                  false,
						ShowDestinationSelector: false,
						Directions:              []string{"DIR_1"},
					},
				},
				Directions: []Direction{
					{
						ID:        "DIR_1",
						ShortName: "Inbound",
						Name:      "Inbound to Downtown",
						UseForUI:  true,
						Stops:     []string{"3860"},
					},
				},
				Paths: []Path{
					{
						ID: "PATH_1",
						Points: []PathPoint{
							{Lat: 37.7936799, Lon: -122.39637},
							{Lat: 37.7935899, Lon: -122.39746},
						},
					},
				},
				Timestamp: "2025-04-26T10:31:08Z",
			}, nil
		},
		GetPredictionsFunc: func(ctx context.Context, routeID, stopID string) ([]Prediction, error) {
			if routeID == "" {
				return nil, ErrRouteIDRequired
			}

			if stopID == "" {
				return nil, ErrStopIDRequired
			}

			return []Prediction{
				{
					VehicleID:       "51",
					Minutes:         9,
					Direction:       "Market & California",
					DestinationName: "Market & California",
					Timestamp:       time.Now().Add(9 * time.Minute),
					VehicleType:     "Cable Car_CABLECAR",
					IsDeparture:     true,
				},
				{
					VehicleID:       "59",
					Minutes:         19,
					Direction:       "Market & California",
					DestinationName: "Market & California",
					Timestamp:       time.Now().Add(19 * time.Minute),
					VehicleType:     "Cable Car_CABLECAR",
					IsDeparture:     true,
				},
			}, nil
		},
		ClearCacheFunc: func() {
			// Do nothing in the mock
		},
		EnableCacheFunc: func() {
			// Do nothing in the mock
		},
		DisableCacheFunc: func() {
			// Do nothing in the mock
		},
	}
}

// GetAllRoutes calls the mock implementation
func (m *MockClient) GetAllRoutes(ctx context.Context) ([]RouteInfo, error) {
	return m.GetAllRoutesFunc(ctx)
}

// GetRouteDetails calls the mock implementation
func (m *MockClient) GetRouteDetails(ctx context.Context, routeID string) (*RouteDetails, error) {
	return m.GetRouteDetailsFunc(ctx, routeID)
}

// GetPredictions calls the mock implementation
func (m *MockClient) GetPredictions(ctx context.Context, routeID, stopID string) ([]Prediction, error) {
	return m.GetPredictionsFunc(ctx, routeID, stopID)
}

// ClearCache calls the mock implementation
func (m *MockClient) ClearCache() {
	m.ClearCacheFunc()
}

// EnableCache calls the mock implementation
func (m *MockClient) EnableCache() {
	m.EnableCacheFunc()
}

// DisableCache calls the mock implementation
func (m *MockClient) DisableCache() {
	m.DisableCacheFunc()
}
