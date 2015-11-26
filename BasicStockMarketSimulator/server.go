package main

import (
    "fmt"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "os"
	"strings"
    "encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"strconv"
	"errors"
)
	
type StockQuote struct {
    List struct{
	Meta struct{
		Count int	`json:"count"`
   		Start int 	`json:"start"`
		Type string `json:"type"`
	} `json:"meta"`
	Resources []struct{
		Resource struct{
			Classname string `json:"classname"`
			Fields struct{
				Name string    `json:"name"`
				Price string   `json:"price"`
				Symbol string  `json:"symbol"`
				Ts string      `json:"ts"`
				Type string    `json:"type"`
				UTCtime string `json:"utctime"`
				Volume string  `json:"volume"`
			}`json:"fields"`
		}`json:"resource"`
	}`json:"resources"`
    }`json:"list"`
}

type BuyStockRequest struct {
	Budget float64                  `json:"budget"`
	StockSymbolAndPercentage string `json:"stockSymbolAndPercentage"`
}

type BuyStockResponse struct{
	TradeId uint32         `json:"tradeid"`
	Stocks string          `json:"stocks"`
	UnvestedAmount float64 `json:"unvestedAmount"`
}

type PortfolioRequest struct {
	TradeId uint32 `json:"tradeid"`
}

type PortfolioResponse struct {
	Stocks string              `json:"stocks"`
	CurrentMarketValue float64 `json:"currentMarketValue"`
	UnvestedAmount float64     `json:"unvestedAmount"`
}

type Arith int
var stock StockQuote
var trans BuyStockResponse

func (t *Arith) BuyStock(buyStockRequest *BuyStockRequest, buyStockResponse *BuyStockResponse) error {
    
	buyStockRequest.StockSymbolAndPercentage = strings.Replace(buyStockRequest.StockSymbolAndPercentage,":",",",strings.Count(buyStockRequest.StockSymbolAndPercentage,":"))
	buyStockRequest.StockSymbolAndPercentage = strings.Replace(buyStockRequest.StockSymbolAndPercentage,"%","",strings.Count(buyStockRequest.StockSymbolAndPercentage,"%"))
	list := strings.Split(buyStockRequest.StockSymbolAndPercentage,",")
	
	stock :=""
	percent :=""
	var percentTotal float64
	var percnt float64
	
	
	for i:=0;i<len(list);i++	{
		if i%2==0	{
			stock=stock + list[i] + ","
		} 
		if i%2!=0	{
			percent=percent + list[i] + ","
			percnt, _ = strconv.ParseFloat(list[i],64)
			percentTotal += percnt
		}	
	}
	if(percentTotal<=100){
		err := processStocks(stock,percent,buyStockRequest.Budget)
		*buyStockResponse = trans
		
		if err != nil {
			return errors.New("Total quote price exceeds balance")
		}
		
	}else {
		err := errors.New("Invalid percentage of budget")
		return err
	}
	
	
	return nil
}

//timeout constant
const(
	timeout = time.Duration(time.Second*100)
)

func getQuotes(str string)	{

    client := http.Client{Timeout: timeout}
    url := fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json", str)
    res, err := client.Get(url)
    if err != nil {
        fmt.Errorf("Stocks cannot access yahoo finance API: %v", err)
    }
    defer res.Body.Close()
	
    content, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Errorf("Stocks cannot read json body: %v", err)
    }
	
    err = json.Unmarshal(content, &stock)
    if err != nil {
        fmt.Errorf("Stocks cannot parse json data: %v", err)
    }
}

func parseQuotes() (prices string) {
	
	var priceStr  = ""
	count := stock.List.Meta.Count
	for i:=0;i<count;i++ {
		priceStr = priceStr + stock.List.Resources[i].Resource.Fields.Price + ","
	}
	return priceStr
}

func processStocks(stockStr string, percntStr string, balance float64) error {
	
	//call to request stock quotes from Yahoo
	getQuotes(stockStr)
	priceStr:= parseQuotes()
		
	prices := strings.Split(priceStr,",")
	percnts:= strings.Split(percntStr,",")
	stocks := strings.Split(stockStr,",")
	
	var prc float64
	var prcnt float64
	var qty int
	var qtyStr = ""
	var total float64
	var unvestedAmt float64
	stockCount := len(percnts)
	
	for i:=0; i<stockCount-1; i++ {
		prc, _ = strconv.ParseFloat(prices[i],64)
		prcnt, _ = strconv.ParseFloat(percnts[i],64)
		qty = int((balance*prcnt)/(100.00*prc))
		total = total + (float64(qty)*prc)
		qtyStr = qtyStr + strconv.Itoa(qty) +","
	} 
	
	stockQty := strings.Split(qtyStr,",")
	transDetails := ""

	if(total<balance){
		for i:=0; i<stockCount-1; i++ {
		
			prc,_ = strconv.ParseFloat(prices[i],64)
			transDetails= transDetails + stocks[i] + ":" + stockQty[i] +":$" + prices[i]
			//fmt.Print(transDetails)
			if(i!=stockCount-2){
				//fmt.Printf(",")
				transDetails += ","
				
			}
		}
		


	}else{
		err:= errors.New("Exceeds balance")
		return err
	}
		unvestedAmt = balance-total
		t := time.Unix(0, 11/11/2011)
		trans.TradeId = uint32(time.Since(t))
		trans.Stocks = transDetails
		trans.UnvestedAmount = unvestedAmt
	return nil
}

func (t *Arith) CheckPortfolio(portfolioRequest *PortfolioRequest, portfolioResponse *PortfolioResponse) error {
	
	if(portfolioRequest.TradeId == trans.TradeId){
		trans.Stocks = strings.Replace(trans.Stocks,"$","",strings.Count(trans.Stocks,"$"))
		trans.Stocks = strings.Replace(trans.Stocks,":",",",strings.Count(trans.Stocks,":"))
		list := strings.Split(trans.Stocks,",")
			
		stockStr:=""
		qtyStr:=""
		oldPriceStr:=""
			
		for i:=0;i<len(list);i++	{
			if i%3==0	{
				stockStr = stockStr + list[i] + ","
			} 
			if i%3==1	{
				qtyStr = qtyStr + list[i] + ","
			}	
			if i%3==2   {
				oldPriceStr = oldPriceStr + list[i] + ","
			}
		}
		
		//get current price using getQuotes
		getQuotes(stockStr)
		newPriceStr:= parseQuotes()	
				
		newPriceArry := strings.Split(newPriceStr,",")
		oldPriceArry := strings.Split(oldPriceStr,",")
		symbolsArry := strings.Split(stockStr,",")
		qtyArry := strings.Split(qtyStr,",")
		
		var newPrc float64
		var oldPrc float64
		var qty int64
		var total float64
		var stocksStr string
		
		stockCount := len(qtyArry)
			
		for i:=0; i<stockCount-1; i++ {
			newPrc, _ = strconv.ParseFloat(newPriceArry[i],64)
			oldPrc, _ = strconv.ParseFloat(oldPriceArry[i],64)
			qty, _ = strconv.ParseInt(qtyArry[i], 10, 64)
			total = total + (float64(qty)*newPrc)
			if newPrc > oldPrc {
				stocksStr = stocksStr + symbolsArry[i] + ":" + qtyArry[i] + ":+$" + newPriceArry[i]
			} else if newPrc < oldPrc {
				stocksStr = stocksStr + symbolsArry[i] + ":" + qtyArry[i] + ":-$" + newPriceArry[i]
			}else {
				stocksStr = stocksStr + symbolsArry[i] + ":" + qtyArry[i] + ":$" + newPriceArry[i]
			}
			if(i!=stockCount-2){
				stocksStr += ","
			}
		}

		portfolioResponse.Stocks = stocksStr
		portfolioResponse.CurrentMarketValue = total
		portfolioResponse.UnvestedAmount = trans.UnvestedAmount
		
	}
	return nil
}

func main() {

    arith := new(Arith)
    rpc.Register(arith)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }

}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}