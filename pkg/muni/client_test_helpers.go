package muni

import (
	"net/http"
	"net/http/httptest"
)

// mockServer creates a test server that returns the given response for all requests
func mockServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

var mockRoutesResponse = `[
	{
		"id": "N",
		"rev": 1,
		"title": "N-Judah",
		"description": "N-Judah Line",
		"color": "003399",
		"textColor": "FFFFFF",
		"hidden": false,
		"timestamp": "2024-03-20T12:00:00Z"
	},
	{
		"id": "J",
		"rev": 1,
		"title": "J-Church",
		"description": "J-Church Line",
		"color": "339900",
		"textColor": "FFFFFF",
		"hidden": false,
		"timestamp": "2024-03-20T12:00:00Z"
	}
]`

var mockRouteDetailsResponse = `{
	"id": "N",
	"rev": 1,
	"title": "N-Judah",
	"description": "N-Judah Line",
	"color": "003399",
	"textColor": "FFFFFF",
	"hidden": false,
	"boundingBox": {
		"latMin": 37.7601,
		"latMax": 37.7749,
		"lonMin": -122.5089,
		"lonMax": -122.3894
	},
	"stops": [
		{
			"id": "1234",
			"lat": 37.7749,
			"lon": -122.4194,
			"name": "Ocean Beach",
			"hidden": false,
			"showDestinationSelector": true,
			"directions": ["Inbound", "Outbound"]
		}
	],
	"directions": [
		{
			"id": "IB",
			"shortName": "IB",
			"name": "Inbound",
			"useForUi": true,
			"stops": ["1234"]
		}
	],
	"paths": [
		{
			"id": "1",
			"points": [
				{
					"lat": 37.7749,
					"lon": -122.4194
				}
			]
		}
	],
	"timestamp": "2024-03-20T12:00:00Z"
}`

var mockPredictionsResponse = `[{
	"serverTimestamp": 1710936000,
	"nxbs2RedirectUrl": "",
	"agency": {
		"rev": 1,
		"id": "sfmta-cis",
		"name": "San Francisco Municipal Transportation Agency",
		"shortName": "SFMTA"
	},
	"route": {
		"id": "N",
		"title": "N-Judah",
		"description": "N-Judah Line",
		"color": "003399",
		"textColor": "FFFFFF",
		"hidden": false
	},
	"stop": {
		"id": "1234",
		"lat": 37.7749,
		"lon": -122.4194,
		"name": "Ocean Beach",
		"hidden": false,
		"showDestinationSelector": true,
		"route": "N"
	},
	"values": [
		{
			"timestamp": 1710936000,
			"minutes": 5,
			"affectedByLayover": false,
			"isDeparture": false,
			"occupancyStatus": 1,
			"occupancyDescription": "Many Seats Available",
			"vehiclesInConsist": 1,
			"linkedVehicleIds": "",
			"vehicleId": "1234",
			"vehicleType": "LRV4",
			"direction": {
				"id": "IB",
				"name": "Inbound",
				"destinationName": "Downtown"
			},
			"tripId": "1234",
			"delay": 0,
			"predUsingNavigationTm": false,
			"departure": false
		}
	]
}]`
