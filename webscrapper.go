package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup
var jsonResponces = make(chan string)

func scrapper(url string) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			jsonResponces <- string(body)
		}
	}
}

func main() {

	urls := []string{
		"http://www.reddit.com/r/aww.json",
		"http://www.reddit.com/r/funny.json",
		"http://www.reddit.com/r/programming.json",
	}

	// Starting scrapping
	fmt.Println("main")
	go func() {
		for {
			for _, url := range urls {
				wg.Add(1)
				go scrapper(url)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Printing result
	wg.Add(1)
	go func() {
		for response := range jsonResponces {
			fmt.Println(response)
		}
	}()
	// Waiting for all program to complete
	wg.Wait()
	fmt.Println(" Program Shutdown...")
}
