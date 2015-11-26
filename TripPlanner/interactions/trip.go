package interactions

import (
	"gopkg.in/mgo.v2/bson"
    "fmt"
	"components"
    "encoding/json"
    "net/http"
    "time"
   	"github.com/julienschmidt/httprouter"
)

const(
    timeout = time.Duration(time.Second*100)
)

func (sess MongodbHandler) GetTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	tripDetails := components.TripDetails{}
	id := bson.ObjectIdHex(p.ByName("trip_id"))
	err := sess.session.DB("test_db").C("trips").FindId(id).One(&tripDetails)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    } 	
	if(tripDetails.ID == ""){
		fmt.Fprintf(rw,"Invalid LocationID")
	}else{
		tripDetails.NextDestinationLocationID=""
		tripDetails.UberWaitTime= ""
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
    	rw.WriteHeader(200)
    	json.NewEncoder(rw).Encode(tripDetails)
	}
}


func (sess MongodbHandler) AddTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	tripRequest := components.TripRequest{}
	locationDetails := components.Location{}
	tripDetails := components.TripDetails{}
	locations := make(map[string]components.Location)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&tripRequest)
	if(err!=nil)    {
		fmt.Errorf("Error in decoding the Input JSON: %v", err)
	}
	
	url := fmt.Sprintf("http://localhost:8080/locations/%s",tripRequest.StartingFromLocationID)
	client := http.Client{Timeout: timeout}

    res, err := client.Get(url)
    if err != nil {
        fmt.Errorf("Cannot read localhost LocationsAPI: %v", err)
    }
    defer res.Body.Close()
    decoder = json.NewDecoder(res.Body)
	
    err = decoder.Decode(&locationDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Location JSON: %v", err)
    }
    
	//Push the startingID location details to map
	locations[tripRequest.StartingFromLocationID] = locationDetails
	
	//Push the rest of the IDs' location details to map by iteration
	for _, value := range tripRequest.LocationIds{
		url = fmt.Sprintf("http://localhost:8080/locations/%s",value)
		client = http.Client{Timeout: timeout}

		res, err = client.Get(url)
		if err != nil {
			fmt.Errorf("Cannot read localhost LocationsAPI: %v", err)
		}
		defer res.Body.Close()
		decoder = json.NewDecoder(res.Body)
		
		err = decoder.Decode(&locationDetails)
		if(err!=nil)    {
			fmt.Errorf("Error in decoding the Location JSON: %v", err)
		}
		locations[value] = locationDetails
	}
	startID := tripRequest.StartingFromLocationID
	startLat := locations[startID].Coordinate.Lat
	originLat := startLat
	startLng := locations[startID].Coordinate.Lng
	originLng := startLng
	nextID := tripRequest.StartingFromLocationID
	var lowPrice int
	var duration int
	var distance float64
	minPrice := 10000 
	minduration := 0
	mindistance := 0.0
	totalCost := 0
	totalUberDuration := 0
	totalDistance := 0.0
	locationOrder := 0
	fmt.Println(totalDistance)
	for len(locations) > 1 {
		for key, value := range locations{
			if(key != startID){
				
				lowPrice,duration,distance = components.PriceEstimate(startLat,startLng,value.Coordinate.Lat,value.Coordinate.Lng)
				if(lowPrice < minPrice){
					minPrice = lowPrice
					minduration = duration
					mindistance = distance
					nextID = key		
				}	
			}
		}
		totalCost += lowPrice
		totalUberDuration += minduration
		totalDistance += mindistance
		delete(locations,startID)
		startID = nextID
		startLat = locations[startID].Coordinate.Lat
		startLng = locations[startID].Coordinate.Lng
		tripRequest.LocationIds[locationOrder]=nextID
		locationOrder++
		minPrice = 1000000.0
		minduration = 0
		mindistance = 0.0
	}

	//Estimating the price/durtion from last location point to starting point for a round-trip
	lowPrice,duration,distance = components.PriceEstimate(originLat,originLng,locations[nextID].Coordinate.Lat,locations[nextID].Coordinate.Lng)
	totalCost += lowPrice
	totalUberDuration += duration
	totalDistance += distance

	tripDetails.BestRouteLocationIds = tripRequest.LocationIds
	tripDetails.StartingFromLocationID = tripRequest.StartingFromLocationID
	tripDetails.Status = "planning"
	tripDetails.TotalDistance = totalDistance
	tripDetails.TotalUberDuration = totalUberDuration
	tripDetails.TotalUberCosts = totalCost
	
	tripDetails.ID = bson.NewObjectId()
	err = sess.session.DB("test_db").C("trips").Insert(tripDetails)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
	}
	err = sess.session.DB("test_db").C("trips").FindId(tripDetails.ID).One(&tripDetails)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	}
	
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(tripDetails)
}

func (sess MongodbHandler) UpdateTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	tripDetails := components.TripDetails{}
	locationDetails := components.Location{}
	id := bson.ObjectIdHex(p.ByName("trip_id"))
	err := sess.session.DB("test_db").C("trips").FindId(id).One(&tripDetails)
	if err != nil {
		fmt.Printf("got an error finding a trip %v\n")
	}
	currentLocationID := tripDetails.StartingFromLocationID
	if(tripDetails.NextDestinationLocationID == tripDetails.StartingFromLocationID){
		tripDetails.Status = "trip over"
		tripDetails.NextDestinationLocationID = ""
		tripDetails.StartingFromLocationID = ""
		tripDetails.UberWaitTime = ""
	}else{
				
		if(tripDetails.Status == "requesting"){
			if(len(tripDetails.BestRouteLocationIds)>1){
				currentLocationID = tripDetails.BestRouteLocationIds[0]
				x := tripDetails.BestRouteLocationIds[1:len(tripDetails.BestRouteLocationIds)]
				tripDetails.BestRouteLocationIds = x
				tripDetails.NextDestinationLocationID = tripDetails.BestRouteLocationIds[0]
			}else{
				tripDetails.BestRouteLocationIds = nil
				currentLocationID = tripDetails.NextDestinationLocationID
				tripDetails.NextDestinationLocationID = tripDetails.StartingFromLocationID
			}
		}else if(tripDetails.Status == "planning"){
			tripDetails.NextDestinationLocationID = tripDetails.BestRouteLocationIds[0]	
			tripDetails.Status = "requesting"
		}

		url := fmt.Sprintf("http://localhost:8080/locations/%s",currentLocationID)
		client := http.Client{Timeout: timeout}

	    res, err := client.Get(url)
	    if err != nil {
	        fmt.Errorf("Cannot read localhost LocationsAPI: %v", err)
	    }
	    defer res.Body.Close()
	    decoder := json.NewDecoder(res.Body)
		
	    err = decoder.Decode(&locationDetails)
	    if(err!=nil)    {
	        fmt.Errorf("Error in decoding the Location JSON: %v", err)
	    }
		startLat:=locationDetails.Coordinate.Lat
		startLng:=locationDetails.Coordinate.Lng

		url = fmt.Sprintf("http://localhost:8080/locations/%s",tripDetails.NextDestinationLocationID)
		client = http.Client{Timeout: timeout}

	    res, err = client.Get(url)
	    if err != nil {
	        fmt.Errorf("Cannot read localhost LocationsAPI: %v", err)
	    }
	    defer res.Body.Close()
	    decoder = json.NewDecoder(res.Body)
		
	    err = decoder.Decode(&locationDetails)
	    if(err!=nil)    {
	        fmt.Errorf("Error in decoding the Location JSON: %v", err)
	    }
		tripDetails.UberWaitTime = components.UberRideRequest(startLat,startLng,locationDetails.Coordinate.Lat,locationDetails.Coordinate.Lng)
		

		if(len(tripDetails.BestRouteLocationIds)==0){
			tripDetails.NextDestinationLocationID = tripDetails.StartingFromLocationID
		}
		///////////////////////////////////////
	}

	//update the request in database
	err = sess.session.DB("test_db").C("trips").UpdateId(id,tripDetails)
  	if err != nil {
    	fmt.Printf("got an error updating a doc %v\n")
    } 
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(tripDetails)
}





