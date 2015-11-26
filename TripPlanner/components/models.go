package components

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

const(
    timeout = time.Duration(time.Second*100)
)

type Location struct {
    ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Address    string `json:"address"`
    City       string `json:"city"`
    Name  string `json:"name"`
    State string `json:"state"`
    Zip   string `json:"zip"`
    Coordinate struct {
        	Lat float64 `json:"lat"`
        	Lng float64 `json:"lng"`
    } `json:"coordinate"`
} 

type GoogleLocation struct {
	    Results []struct {
	        AddressComponents []struct {
	            LongName  string   `json:"long_name"`
	            ShortName string   `json:"short_name"`
	            Types     []string `json:"types"`
	        } `json:"address_components"`
	        FormattedAddress string `json:"formatted_address"`
	        Geometry         struct {
	            Location struct {
	                Lat float64 `json:"lat"`
	                Lng float64 `json:"lng"`
	            } `json:"location"`
	            LocationType string `json:"location_type"`
	            Viewport     struct {
	                Northeast struct {
	                    Lat float64 `json:"lat"`
	                    Lng float64 `json:"lng"`
	                } `json:"northeast"`
	                Southwest struct {
	                    Lat float64 `json:"lat"`
	                    Lng float64 `json:"lng"`
	                } `json:"southwest"`
	            } `json:"viewport"`
	        } `json:"geometry"`
	        PlaceID string   `json:"place_id"`
	        Types   []string `json:"types"`
	    } `json:"results"`
	    Status string `json:"status"`
}

type TripRequest struct {
	LocationIds            []string `json:"location_ids"`
	StartingFromLocationID string   `json:"starting_from_location_id"`
}

type TripDetails struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Status                 string   `json:"status"`	
	StartingFromLocationID string   `json:"starting_from_location_id,omitempty"`
	NextDestinationLocationID string `json:"next_destination_location_id,omitempty"`
	BestRouteLocationIds   []string `json:"best_route_location_ids"`
	TotalUberCosts         int      `json:"total_uber_costs"`
	TotalUberDuration      int      `json:"total_uber_duration"`
	TotalDistance          float64  `json:"total_distance"`
	UberWaitTime		   string 		`json:"uber_wait_time_eta,omitempty"`
}

type RideRequest struct {
	ProductID      string `json:"product_id"`
	StartLatitude  string `json:"start_latitude"`
	StartLongitude string `json:"start_longitude"`
	EndLatitude    string `json:"end_latitude"`
	EndLongitude   string `json:"end_longitude"`
}

type RideResponse struct {
	Driver          interface{} `json:"driver"`
	Eta             int         `json:"eta"`
	Location        interface{} `json:"location"`
	RequestID       string      `json:"request_id"`
	Status          string      `json:"status"`
	SurgeMultiplier int         `json:"surge_multiplier"`
	Vehicle         interface{} `json:"vehicle"`
}

type UberProducts struct {
	Products []struct {
		Capacity     int    `json:"capacity"`
		Description  string `json:"description"`
		DisplayName  string `json:"display_name"`
		Image        string `json:"image"`
		PriceDetails struct {
			Base            float64 `json:"base"`
			CancellationFee int     `json:"cancellation_fee"`
			CostPerDistance float64 `json:"cost_per_distance"`
			CostPerMinute   float64 `json:"cost_per_minute"`
			CurrencyCode    string  `json:"currency_code"`
			DistanceUnit    string  `json:"distance_unit"`
			Minimum         float64 `json:"minimum"`
			ServiceFees     []struct {
				Fee  float64 `json:"fee"`
				Name string  `json:"name"`
			} `json:"service_fees"`
		} `json:"price_details"`
		ProductID string `json:"product_id"`
	} `json:"products"`
}

type UberPriceEstimates struct {
	Prices []struct {
		CurrencyCode    string  `json:"currency_code"`
		DisplayName     string  `json:"display_name"`
		Distance        float64 `json:"distance"`
		Duration        int     `json:"duration"`
		Estimate        string  `json:"estimate"`
		HighEstimate    int     `json:"high_estimate"`
		LowEstimate     int     `json:"low_estimate"`
		ProductID       string  `json:"product_id"`
		SurgeMultiplier int     `json:"surge_multiplier"`
	} `json:"prices"`
}