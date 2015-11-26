package main

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"encoding/json"
)

type DataSet struct {
    Key    int `json:"key"`
    Value  string `json:"value"`
}

type DataStore struct {
	DataPoints []DataSet `json:"data"`
}

//var dataMap map[string]string
var dataMap = make(map[int]DataSet)

func ReturnAllKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	dataStore := make([]DataSet, len(dataMap))
	index := 0
	for _,value := range  dataMap {
		dataStore[index] = value
		index++
	}
	json.NewEncoder(rw).Encode(dataStore)
}

func ReturnKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	key,_ := strconv.Atoi(p.ByName("key_id"))
	if dataMap[key].Value != "" {
			json.NewEncoder(rw).Encode(dataMap[key])
		}else {
			fmt.Fprintf(rw,"Invalid Key")
		}
	
}

func SaveKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
	var dataSet DataSet
	dataSet.Key,_ = strconv.Atoi(p.ByName("key_id"))
	dataSet.Value = p.ByName("value")
	dataMap[dataSet.Key] = dataSet
	rw.WriteHeader(200)
}

func main() {

	mux := httprouter.New()
    	
	mux.GET("/keys", ReturnAllKeyValue)
	mux.GET("/keys/:key_id", ReturnKeyValue)
    mux.PUT("/keys/:key_id/:value",SaveKeyValue)
	server := http.Server{
            Addr:        "0.0.0.0:3000",
            Handler: mux,
    }
    
    server.ListenAndServe()

}

