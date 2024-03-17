package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	s "slices"
	"strconv"
	"strings"
	"sync"
	"time"

	concurrent_queue "github.com/fahimfaisaal/exp-go-concurrency/concurrent-queue"
	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
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

func runTasks() {
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

func runServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// http://localhost:4000/cpu-task/compute?n=10000000000
	r.Get("/cpu-task/compute", func(w http.ResponseWriter, r *http.Request) {
		n, err := strconv.ParseUint(r.URL.Query().Get("n"), 10, 64)

		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		task(Task{n, "CPU-Task"})

		w.Write([]byte("Done!"))
	})

	// http://localhost:4000/cpu-tasks/compute?c=5,numbers=10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000,10000000000
	r.Get("/cpu-tasks/compute", func(w http.ResponseWriter, r *http.Request) {
		n := strings.Split(r.URL.Query().Get("numbers"), ",")
		c64, err := strconv.ParseUint(r.URL.Query().Get("c"), 10, 32)
		concurrency := uint(c64)

		if err != nil {
			concurrency = uint(math.Min(float64(len(n)/2), 30))
		}

		q := concurrent_queue.New[Task](concurrency, task)

		for index, num := range n {
			num, err := strconv.ParseUint(num, 10, 64)

			if err != nil {
				q.Close()
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			q.Add(Task{num, fmt.Sprintf("CPU-Task-%d", index+1)})
		}

		q.Close()
		w.Write([]byte("Done!"))
	})

	r.Get("/io-task/wait", func(w http.ResponseWriter, r *http.Request) {
		delay, err := strconv.ParseUint(r.URL.Query().Get("delay"), 10, 32)

		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		time.Sleep(time.Duration(delay) * time.Second)

		w.Write([]byte("Done!"))
	})

	logger := log.New(os.Stdout, "[SERVER] ", log.Ldate|log.Ltime)

	server := &http.Server{
		Addr:     ":4000",
		Handler:  r,
		ErrorLog: logger,
	}

	logger.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

func main() {
	if s.Contains(os.Args[1:], "server") || os.Getenv("TYPE") == "server" {
		runServer()
	} else {
		runTasks()
	}
}
