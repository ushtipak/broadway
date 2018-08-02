package main

import (
	"log"
	"time"
	"net/http"
	"io"
	"os"
	"io/ioutil"
)

var (
	port               = "7835"
	logFd              = "output.log"
	votes              []string
	logDebug, logError *log.Logger
)

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

	logDebug.Printf("brodway listen on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func voteHandler(_ http.ResponseWriter, r *http.Request) {
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError.Println(err)
	}
	logDebug.Printf("%s", resp)
}