package main

import (
	"crypto/md5"
	"fmt"
//	"bytes"
	"encoding/json"
//	"encoding/binary"
	"net/http"
    "time"
    "strconv"

)

type DataSet struct {
    Key    string `json:"key"`
    Value  string `json:"value"`
}

const(
    timeout = time.Duration(time.Second*100)
)

func main() {
	
	var mod uint8
	dataSet := DataSet{}

	ch := 'a'
	for key := 1; key<=10; key++ {
		dataSet.Key   = strconv.Itoa(key)
		dataSet.Value = string(ch)
		dataSetJSON,_ := json.Marshal(dataSet)
		hashValue := md5.Sum(dataSetJSON)
		
		for i:=0;i<16;i++{
			mod+=hashValue[i]
		}
		cache := mod%3
		fmt.Printf("%s - %s : %d\n",dataSet.Key,dataSet.Value,cache)
		    url := fmt.Sprintf("http://localhost:300%d/keys/%s/%s",cache,dataSet.Key,dataSet.Value)
		    fmt.Println(url)
			client := http.Client{Timeout: timeout}
			req, err := http.NewRequest("PUT", url, nil)
			if err != nil {
		        fmt.Errorf("Error in making PUT request")
		    }
			res, err := client.Do(req)
			defer res.Body.Close()
		    if err != nil {
		        fmt.Errorf("Error in calling PUT request")
		    }
		ch++  
		mod = 0  
	}
}