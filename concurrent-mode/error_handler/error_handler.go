package main

import (
	"fmt"
	"net/http"
)

func main() {

}

func handleErr() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("response: %v\n", response.Status)
	}
}

type Result struct {
	Error    error
	Response *http.Response
}

func handleErr2() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		responses := make(chan Result)
		go func() {
			defer close(responses)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case responses <- result:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("err: %v\n", result.Error)
			continue
		}
		fmt.Printf("response: %v\n", result.Response.Status)
	}
}
