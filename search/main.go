package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/crywolf/goexperiments/search/fsearch"
)

var mode *string

type response struct {
	Results []fsearch.Result
	Elapsed time.Duration
}

func main() {
	port := flag.String("port", "8080", "defines port")
	mode = flag.String("mode", "serial", "search mode [serial, parallel, timeout]")
	flag.Parse()

	if !isModeValid(*mode) {
		fmt.Println("incorrect value of '-mode' flag")
		os.Exit(1)
	}

	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/health", handleHealth)
	log.Printf("serving on http://localhost:%s/search\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "missing 'q' parameter", http.StatusBadRequest)
		return
	}

	var results []fsearch.Result
	var err error

	start := time.Now()

	switch *mode {
	case "serial":
		results, err = fsearch.Search(query)
	case "parallel":
		results, err = fsearch.ParallelSearch(query)
	case "timeout":
		results, err = fsearch.TimeoutSearch(query, 80*time.Millisecond)
	default:
		log.Fatal("incorrect '-mode' parameter")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	elapsed := time.Since(start)

	res := response{results, elapsed}

	switch req.FormValue("output") {
	case "json":
		err = json.NewEncoder(w).Encode(res)
	case "prettyjson":
		var b []byte
		b, err = json.MarshalIndent(res, "", " ")
		_, err = w.Write(b)
	default:
		fmt.Fprintf(w, "%+v", res)
	}

	if err != nil {
		msg := fmt.Sprintf("could not encode search result: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func handleHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func isModeValid(mode string) bool {
	switch mode {
	case
		"serial",
		"parallel",
		"timeout":
		return true
	}
	return false
}
