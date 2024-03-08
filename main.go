package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	concurrent_queue "github.com/fahimfaisaal/exp-go-concurrency/concurrent-queue"
)

type Task struct {
	length   uint64
	taskName string
}

var pid int = os.Getpid()

func task(task Task) {
	fmt.Printf("Started %s, PID: %d\n", task.taskName, pid)
	sum := uint64(0)

	for i := range task.length {
		sum += i
		// fmt.Printf("Working on %s\n", task.taskName)
	}

	fmt.Printf("Done %s PID: %d\n", task.taskName, pid)
}

func ioTask(endpoint string) {
	fmt.Printf("Checking the status of %s\n", endpoint)
	res, err := http.Get(endpoint)

	if err != nil {
		println("Oh no! An error occurred!")
		return
	}

	fmt.Printf("The status code is: %d\n", res.StatusCode)
}

func main() {
	start := time.Now()
	defer func() {
		end := time.Now()
		fmt.Printf("Execution time: %s\n", end.Sub(start))
	}()

	q := concurrent_queue.New[Task](3, task)
	q2 := concurrent_queue.New[string](3, ioTask)
	defer q2.Close()
	defer q.Close()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for i := range 6 {
			q.Add(Task{1e10, fmt.Sprintf("Task-%d", i+1)})
		}
		fmt.Println("All tasks been added")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
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
			q2.Add(link)
		}

		wg.Done()
	}()

	wg.Wait()
}
