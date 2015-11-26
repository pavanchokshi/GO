package components

import (
	"fmt"
    "encoding/json"
    "net/http"
    "encoding/binary"
    "bytes"
    "strconv"
)



func GetUberProductID(latitude float64, longitude float64) string {

	uberDetails := UberProducts{}
    url := fmt.Sprintf("%s/products?latitude=%f&longitude=%f",UBER_URL,latitude,longitude)
	client := http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization",UBER_SERVER_TOKEN)
	res, err := client.Do(req)
	defer res.Body.Close()
    if err != nil {
        fmt.Errorf("Cannot read UBER API: %v", err)
    }
	decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&uberDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }
    return uberDetails.Products[0].ProductID
}

//returns low_estimate and high_estimate
func PriceEstimate(start_latitude float64, start_longitude float64, end_latitude float64, end_longitude float64)(int,int,float64) {

	uberPriceEstimates := UberPriceEstimates{}
	url := fmt.Sprintf("%s/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f",UBER_URL,start_latitude,start_longitude,end_latitude,end_longitude)
	client := http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization",UBER_SERVER_TOKEN)
	res, err := client.Do(req)
	defer res.Body.Close()
    if err != nil {
        fmt.Errorf("Cannot read UBER API: %v", err)
    }
	decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&uberPriceEstimates)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }

    return uberPriceEstimates.Prices[0].LowEstimate, uberPriceEstimates.Prices[0].Duration, uberPriceEstimates.Prices[0].Distance
}

func UberRideRequest(start_latitude float64, start_longitude float64, end_latitude float64, end_longitude float64) string {
	
	rideRequest := RideRequest{}
	rideResponse := RideResponse{}
	rideRequest.ProductID = GetUberProductID(start_latitude,start_longitude)
	rideRequest.StartLatitude = fmt.Sprintf("%.6f",start_latitude)
	rideRequest.StartLongitude = fmt.Sprintf("%.6f",start_longitude)
	rideRequest.EndLatitude = fmt.Sprintf("%.6f",end_latitude)
	rideRequest.EndLongitude = fmt.Sprintf("%.6f",end_longitude)

	url := fmt.Sprintf("%s/requests",UBER_URL)
	client := http.Client{Timeout: timeout}
	b, err := json.Marshal(rideRequest)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &b)
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization",UBER_ACCESS_TOKEN)
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
        fmt.Errorf("Error in UBER ride request API: %v", err)
    }

	decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&rideResponse)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the UBER_RideRequest JSON: %v", err)
    }
    eta := strconv.Itoa(rideResponse.Eta)
    fmt.Printf("StartLatitude:%s - StartLongitude%s\n",start_latitude,start_longitude)
    fmt.Printf("StopLatitude:%s - StopLongitude%s\n",end_latitude,end_longitude)
    fmt.Printf("Estimated Wait Time:%s\n",eta)
    fmt.Println("--------------------------------------------------------------")
    return eta

}



