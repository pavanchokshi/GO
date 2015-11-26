package main

import (
	"fmt"
    "log"
    "net/rpc/jsonrpc"
    "os"
	"encoding/json"
	"strconv"	
)

type BuyStockRequest struct {
	Budget float64    				`json:"budget"`
	StockSymbolAndPercentage string `json:"stockSymbolAndPercentage"`
}

type BuyStockResponse struct{
	TradeId uint32 		   `json:"tradeid"`
	Stocks string 		   `json:"stocks"`
	UnvestedAmount float64 `json:"unvestedAmount"`
}

type PortfolioRequest struct {
	TradeId uint32 `json:"tradeid"`
}

type PortfolioResponse struct {
	Stocks string 				`json:"stocks"`
	CurrentMarketValue float64  `json:"currentMarketValue"`
	UnvestedAmount float64 		`json:"unvestedAmount"`
}

var buyStockRequest BuyStockRequest
var buyStockResponse BuyStockResponse
var portfolioRequest PortfolioRequest
var portfolioResponse PortfolioResponse

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: ", os.Args[0], "server:port")
        log.Fatal(1)
    }
	
    service := os.Args[1]
    client, err := jsonrpc.Dial("tcp", service)
    if err != nil {
        log.Fatal("dialing:", err)
    }
	
	if os.Args[2] == "buy"	{
	
		fmt.Println("Buying Stocks.. ")
		content := []byte(os.Args[3])
		err = json.Unmarshal(content, &buyStockRequest)
		err = client.Call("Arith.BuyStock", buyStockRequest, &buyStockResponse)
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Printf("%+v\n",buyStockResponse)
		
	}else if os.Args[2] == "checkPortfolio"	{
		fmt.Println("Checking Portfolio..")
		tradeid,_ := strconv.ParseInt(os.Args[3],10,64)
		portfolioRequest.TradeId = uint32(tradeid)
		err = client.Call("Arith.CheckPortfolio", portfolioRequest , &portfolioResponse)
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Printf("%+v\n",portfolioResponse)
		
	}else	{
		
		fmt.Printf("Invalid Input")
	}
}