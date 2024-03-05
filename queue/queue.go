package queue

import (
	"sync"
)

type Queue[T any] struct {
	concurrency  uint
	channels     []chan T
	channelIndex uint
	wg           *sync.WaitGroup
}

func New[T any](concurrency uint, processor func(T)) *Queue[T] {
	channels := make([]chan T, concurrency)
	wg := sync.WaitGroup{}

	for i := range concurrency {
		channels[i] = make(chan T)

		wg.Add(1)
		go func(channel <-chan T) {
			defer wg.Done()
			for item := range channel {
				processor(item)
			}
		}(channels[i])
	}

	return &Queue[T]{concurrency, channels, 0, &wg}
}

func (q *Queue[T]) getNextChannel() chan<- T {
	channel := q.channels[q.channelIndex]
	q.channelIndex = (q.channelIndex + 1) % q.concurrency
	return channel
}

func (q *Queue[T]) Add(item T) {
	q.getNextChannel() <- item
}

func (q *Queue[T]) Close() {
	defer q.wg.Wait()
	for _, channel := range q.channels {
		close(channel)
	}
}
