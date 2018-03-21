package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	s := http.Server{Addr: fmt.Sprintf(":%s", *port)}
	go func() {
		log.Print(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("shutdown signal received, exiting...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server shutdown: %v", err)
		os.Exit(1)
	}
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
