package main

import (
	"net/http"
	"components"
	"interactions"
   	"github.com/julienschmidt/httprouter"
)

func main() {

	//components.ProductsUber()
	//low,high := components.PriceEstimate(37.340823,-121.898409,37.335426,-121.884966)
	mongoHandler := interactions.CreateMongodbHandler(components.ConnectMongo())
	mux := httprouter.New()
	//handlers for location functions
	mux.GET("/locations/:locationID", mongoHandler.GetLocation)
    mux.POST("/locations", mongoHandler.AddLocation)
    mux.PUT("/locations/:locationID", mongoHandler.UpdateLocation)
	mux.DELETE("/locations/:locationID", mongoHandler.DeleteLocation)
	mux.OPTIONS("/locations", mongoHandler.GetOptions)
	//handlers for trip functions
	mux.GET("/trips/:trip_id", mongoHandler.GetTrip)
	mux.POST("/trips", mongoHandler.AddTrip)
	mux.PUT("/trips/:trip_id/request", mongoHandler.UpdateTrip)
    
	server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
    server.ListenAndServe()
	
}

