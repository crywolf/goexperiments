package fsearch

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Result represents result of the search
type Result struct {
	Title string
	URL   string
}

// SearchFunc simulates searching for one kind of results
type SearchFunc func(query string) Result

// FakeSearch simulates search engine doing search
func FakeSearch(kind, title, URL string) SearchFunc {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result{
			fmt.Sprintf("%s(%q): %s", kind, query, title),
			URL}
	}
}

var (
	// Web is simulated search for web results
	Web = FakeSearch("web", "The Go programming language", "http://golang.org")
	// Image is simulated search for image results
	Image = FakeSearch("image", "The Go gopher", "http://blog.golang.org/gopher/gopher.png")
	// Video is simulated search for video results
	Video = FakeSearch("video", "Concurrency is not parallelism", "https://www.youtube.com/watch?v=cN_DpYBzKso")
)

// Search simulates sequential search
func Search(query string) ([]Result, error) {
	log.Println("serial search")
	results := []Result{
		Web(query),
		Image(query),
		Video(query),
	}
	return results, nil
}

// ParallelSearch simulates parallel search
func ParallelSearch(query string) ([]Result, error) {
	log.Println("parallel search")
	c := make(chan Result, 3)

	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	results := []Result{
		<-c,
		<-c,
		<-c,
	}
	return results, nil
}

// TimeoutSearch simulates parallel search with timeout
func TimeoutSearch(query string, timeout time.Duration) ([]Result, error) {
	log.Println("parallel search with timeout")
	c := make(chan Result, 3)
	timer := time.After(timeout)

	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	var results []Result
Loop:
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timer:
			for j := i; j < 3; j++ {
				results = append(results, Result{"*** timed out ***", "-"})
			}
			break Loop
		}
	}
	return results, nil
}
