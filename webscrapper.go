package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup

// channel for storing responses
var jsonResponses = make(chan string, 3)

// Function to scrape URL
func reqUrl(url string) {
	// Schedule the call to WaitGroup's Done to tell goroutine is completed.
	defer wg.Done()

	res, err := http.Get(url)
	time.Sleep(3 * time.Second)

	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			jsonResponses <- string(body)
		}
	}
}

func main() {
	urls := []string{
		"http://www.reddit.com/r/funny.json",
		"http://www.reddit.com/r/funny.json",
		"http://www.reddit.com/r/programming.json",
	}

	// Status channel to signal process to stop
	status := make(chan bool)

	// Handling SIGTERM and SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("Completing pending operations. Got signal: ", sig)
		time.Sleep(3 * time.Second)
		status <- false
		fmt.Println("Exiting application")
	}()

	status <- true
	// scrapping till gets a stop signal
	for <-status {
		for _, url := range urls {
			// Adding routines to wait for
			wg.Add(1)
			go reqUrl(url)
			time.Sleep(1 * time.Second)
		}
		// Adding routines to wait for
		wg.Add(1)
		go func() {
			for response := range jsonResponses {
				fmt.Println(response)
			}
		}()
		status <- true
	}

	wg.Wait()
}
