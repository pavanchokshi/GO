package main

import(
	"fmt"
	"net/http"
	"time"
	"os"
	"bytes"
	"io"
	"math"
)

var size int64

func main() {

	if len(os.Args) != 3	{	// Three Command line arguments - Client name , Server IP, Server Port
		fmt.Println("Invalid Arguments")
	}	else	{
		serverIP := os.Args[1]
		serverPort := os.Args[2]
		fmt.Print("Download Speed in Bits/sec:")
		fmt.Println(downloadSpeedTest(serverIP, serverPort))
		fmt.Print("Upload Speed in Bits/sec:")
		fmt.Println(uploadSpeedTest(serverIP, serverPort))
	}

}

func downloadSpeedTest(serverIP string, serverPort string)	int64	{	// Calculates the download Speed

	dt1 := time.Now().UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("http://%s:%s/download",serverIP,serverPort)
	res,err := http.Get(url)
	f, err := os.Create("sample") // Creates a dummyfile for speedtest
	defer f.Close()
	io.Copy(f,res.Body)
	if err != nil	{
		fmt.Print(res)
	}
	size = int64(res.ContentLength)
	dt2 := time.Now().UnixNano() / int64(time.Millisecond)
	downloadTime := int64(dt2-dt1)
	milliTosecs := int64(math.Pow(10, 3))
	return (size*milliTosecs*int64(8))/downloadTime // Returns the calculated download speed
}

func uploadSpeedTest(serverIP string, serverPort string)	int64 {		// Calculates the upload Speed
	file, err := os.Open("sample")
	defer os.Remove("sample")
	defer file.Close()
	if err != nil {
		fmt.Print(err)	
	}
	data := make([]byte, size)
	_, err = file.Read(data)
	r := bytes.NewReader(data)
	ut1 := time.Now().UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("http://%s:%s/upload",serverIP,serverPort)
	_, err = http.Post(url,"application/text",r)
	ut2 := time.Now().UnixNano() / int64(time.Millisecond)
	milliTosecs := int64(math.Pow(10, 3))
	uploadTime := int64(ut2-ut1)
	return (size*milliTosecs*int64(8))/uploadTime // Returns the calculated upload speed		
}