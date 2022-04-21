package utils

import (
	"sync"
)

func getWorkPool(threads int, work func(int)) (chan int, *sync.WaitGroup) {
	var ch = make(chan int, 50)
	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			for {
				index, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				work(index)
			}
		}()
	}

	return ch, &wg
}

func EachLimit(length, limit int, work func(int)) {
	threads := limit
	if length < threads {
		threads = length
	}
	ch, wg := getWorkPool(threads, work)

	for i := 0; i < length; i++ {
		ch <- i
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
}
