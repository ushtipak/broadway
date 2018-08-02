package main

import (
	"log"
	"time"
	"net/http"
	"io"
	"os"
	"io/ioutil"
)

var votes []string
var logDebug, logError *log.Logger
var logFd = "output.log"


func init() {
	logFile, err := os.OpenFile(logFd, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logError.Println(err)
	}
	logDebug = log.New(io.MultiWriter(logFile, os.Stdout), "DEBUG ", log.Ldate|log.Ltime)
	logError = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR ", log.Ldate|log.Ltime)
}

func main() {
	votes = append(votes, time.Now().Format(time.RFC3339))

	mux := http.NewServeMux()
	mux.HandleFunc("/vote", voteHandler)

	log.Fatal(http.ListenAndServe(":7835", mux))
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError.Println(err)
	}
	logDebug.Printf("%s", resp)
}
