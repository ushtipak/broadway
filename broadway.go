package main

import (
	"log"
	"net/http"
	"io"
	"os"
	"strings"
	"html/template"
	"fmt"
	"encoding/json"
)

var (
	port               = "7835"
	logFd              = "output.log"
	templateFd         = "/var/www/html/brodown/template.html"
	resultsFd          = "/var/www/html/brodown/index.html"
	voteResult         = VoteResult{Blue: 0, Yellow: 0}
	renderResult       RenderResult
	logDebug, logError *log.Logger
)

type Vote struct {
	Color string `json:"color"`
}

type VoteResult struct {
	Blue   int
	Yellow int
}

type RenderResult struct {
	Blue   string
	Yellow string
}

func checkErr(err error) {
	if err != nil {
		logError.Println(err)
	}
}

func init() {
	logFile, err := os.OpenFile(logFd, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logError.Println(err)
	}
	logDebug = log.New(io.MultiWriter(logFile, os.Stdout), "DEBUG ", log.Ldate|log.Ltime)
	logError = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR ", log.Ldate|log.Ltime)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/vote", voteHandler)

	logDebug.Printf("brodway listen on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func voteHandler(_ http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)[1]
	if auth != "eC10c2hpcnRicm9kb3duLWF1dGgtdG9rZW4xOmEyNDA4ODY4LTNmMGEtNDViMi1hZDRiLTI1NjUyODk5YTliMg==" {
		logError.Printf("unauthorized")
	} else {
		var vote Vote
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&vote)
		checkErr(err)

		logDebug.Printf("%s", vote)

		if vote.Color == "blue" {
			voteResult.Blue += 1
		}
		if vote.Color == "yellow" {
			voteResult.Yellow += 1
		}
		logDebug.Println(voteResult)

		results, err := os.Create(resultsFd)
		checkErr(err)

		t, err := template.ParseFiles(templateFd)
		checkErr(err)

		renderResult.Blue = fmt.Sprintf("%03d", voteResult.Blue)
		renderResult.Yellow = fmt.Sprintf("%03d", voteResult.Yellow)

		t.Execute(results, renderResult)
	}
}
