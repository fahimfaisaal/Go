package main

import (
	"fmt"
	"net/http"

	"github.com/fahimfaiaal/go/queue"
)

func main() {
	q := queue.New[string](5, func(endpoint string) {
		fmt.Printf("Checking the status of %s\n", endpoint)
		res, err := http.Get(endpoint)

		if err != nil {
			println("Oh no! An error occurred!")
			return
		}

		fmt.Printf("The status code is: %d\n", res.StatusCode)
	})
	defer q.Close()

	links := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.amazon.com",
		"https://www.stackoverflow.com",
		"https://www.discord.com",
		"https://www.httpbin.org",
		"https://www.golang.org",
		"https://www.reddit.com",
		"https://www.x.com",
		"https://www.youtube.com",
	}

	for _, link := range links {
		q.Add(link)
	}
}
