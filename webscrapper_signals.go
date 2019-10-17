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

var jsonResponces = make(chan string)
var quit = make(chan bool)
var terminate = make(chan string)

func scrapper(url string, r int, wg *sync.WaitGroup) {
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
	fmt.Println("Coroutine completed for thread: ", r)
}

func interrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Clenaing up process. Got signal: ", <-sig)
	time.Sleep(3 * time.Second)
	fmt.Println("Quitting the program")
	quit <- true
	close(quit)
}

func main() {

	var wg sync.WaitGroup
	urls := []string{
		"http://ifconfig.co",
		"http://ifconfig.co",
	}

	// Starting signal handler
	go interrupt()

	// Starting scrapping
	fmt.Println("main")
	i := 1
	go func() {
		for {
			select {
			case <-terminate:
				return
			default:
				for _, url := range urls {
					fmt.Println("creating thread ", i)
					// Adding thread to waitgroup
					wg.Add(1)
					go scrapper(url, i, &wg)
					i++
					// Sleeping 1 sec before calling next scrape
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	// Printing result
	// TODO: Have to make it temporary
	go func() {
		for response := range jsonResponces {
			fmt.Println(response)
		}
	}()

	// Waiting for stop signal
	<-quit
	terminate <- "stop"
	close(terminate)

	// Waiting for all program to complete
	wg.Wait()
	fmt.Println(" Program Shutdown...")
}
