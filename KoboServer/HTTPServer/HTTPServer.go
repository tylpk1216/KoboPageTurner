package main

import (
	"io"
	"fmt"
	"net/http"
	"os"
	"time"
)

var count int = 0
var log *os.File

func pingServer(w http.ResponseWriter, r *http.Request) {
    count++
    fmt.Printf("count = %d\n", count)

    io.WriteString(w, "true")

    //s := fmt.Sprintf("count = %d\n", count)
    //log.WriteString(s)
}

func main() {
    log, err := os.Create("HTTPServer.log")
    if err != nil {
        panic(err)
    }

    defer log.Close()

    log.WriteString("Sleep 10 seconds \n")
    time.Sleep(10 * time.Second)
    log.WriteString("Prepare to run Server \n")

    http.HandleFunc("/", pingServer)
    http.ListenAndServe(":80", nil)
}
