package components

import (
    "fmt"
    "strings"
	//"time"
	"encoding/json"
    "net/http"
)

func GetGoogleLocation(address string) (float64, float64)  {

    gmapsDetails := GoogleLocation{}
	address = strings.Replace(address," ","%20",-1)
    url := fmt.Sprintf("http://maps.google.com/maps/api/geocode/json?address=%s",address)
	client := http.Client{Timeout: timeout}

    res, err := client.Get(url)
    if err != nil {
        fmt.Errorf("Cannot read Google API: %v", err)
    }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)

    err = decoder.Decode(&gmapsDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }
    return gmapsDetails.Results[0].Geometry.Location.Lat, gmapsDetails.Results[0].Geometry.Location.Lng
}