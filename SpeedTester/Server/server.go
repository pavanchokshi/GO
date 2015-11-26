package main

import (
	"os"
	"fmt"
	"net/http"
	"io"
)

func main() {
	http.HandleFunc("/upload",Upload) // Handles Upload Speed Test
	http.HandleFunc("/download",Download) // Handles Download Speed Test
	http.ListenAndServe("192.168.2.12:8080", nil) // Supply IP Address and Port
}

// Upload the File to Server
func Upload(rw http.ResponseWriter, r *http.Request) {
    f, err := os.Create("sample") // Creates a dummyfile for speedtest
    defer os.Remove("sample") // Removes the created dummyfile
    defer f.Close()
	io.Copy(f,r.Body)
	if err != nil	{
		fmt.Print(err)
	}	
}

//Serve the file to Client
func Download(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw,r,"./test.txt")
}