package queue

type Queue[T interface{}] struct {
	concurrency  uint
	channels     []chan T
	channelIndex uint
}

func New[T interface{}](concurrency uint, processor func(T)) *Queue[T] {
	channels := make([]chan T, concurrency)

	for i := uint(0); i < concurrency; i++ {
		channels[i] = make(chan T)
	}

	for _, channel := range channels {
		go func(channel <-chan T) {
			for item := range channel {
				processor(item)
			}
		}(channel)
	}

	return &Queue[T]{concurrency, channels, 0}
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
	for _, channel := range q.channels {
		close(channel)
	}
}
