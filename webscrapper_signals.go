package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var wg sync.WaitGroup
var jsonResponces = make(chan string)
var quit = make(chan bool, 2)

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

func interrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig)
	fmt.Println("Clenaing up process. Got signal: ", <-sig)
	time.Sleep(3 * time.Second)
	fmt.Println("Quitting the program")
	quit <- true
}

func main() {

	urls := []string{
		"http://www.reddit.com/r/aww.json",
		"http://www.reddit.com/r/funny.json",
		"http://www.reddit.com/r/programming.json",
	}

	// Starting signal handler
	go interrupt()

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
	go func() {
		for response := range jsonResponces {
			fmt.Println(response)
		}
	}()
	// Waiting to quite the program
	<-quit
	// Waiting for all program to complete
	wg.Wait()
	fmt.Println(" Program Shutdown...")
}
