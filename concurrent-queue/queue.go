package concurrent_queue

import "sync"

type ConcurrentQueue[T any] struct {
	concurrency  uint
	channels     []chan T
	channelIndex uint
	wg           *sync.WaitGroup
}

func New[T any](concurrency uint, processor func(T)) *ConcurrentQueue[T] {
	channels := make([]chan T, concurrency)
	wg := new(sync.WaitGroup)

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

	return &ConcurrentQueue[T]{concurrency, channels, 0, wg}
}

func (q *ConcurrentQueue[T]) getNextChannel() chan<- T {
	channel := q.channels[q.channelIndex]
	q.channelIndex = (q.channelIndex + 1) % q.concurrency
	return channel
}

func (q *ConcurrentQueue[T]) Add(item T) {
	q.getNextChannel() <- item
}

func (q *ConcurrentQueue[T]) Close() {
	defer q.wg.Wait()
	for _, channel := range q.channels {
		close(channel)
	}
}
