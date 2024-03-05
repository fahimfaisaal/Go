package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/fahimfaisaal/go/queue"
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
	q2 := queue.New[float32](5, func(rand float32) {
		time.Sleep(time.Duration(int(rand)) * time.Second)
		fmt.Printf("Wait random time: %f\n", rand)
	})
	defer q.Close()
	defer q2.Close()

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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for _, link := range links {
			fmt.Println("Added link to queue")
			q.Add(link)
		}
		wg.Done()
		fmt.Println("All links added to queue")
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("Add random time")
			q2.Add(rand.Float32() * 10)
		}

		wg.Done()
		fmt.Println("All random number added to queue")
	}()

	wg.Wait()
}
