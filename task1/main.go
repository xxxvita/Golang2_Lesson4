package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var workers = make(chan struct{}, 1)
	var wg sync.WaitGroup
	var countIncrements int = 0

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		workers <- struct{}{}
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			time.Sleep(time.Millisecond)
			_ = <-workers
			countIncrements++
		}(&wg)
	}

	wg.Wait()
	log.Printf("Число увеличений счётчика равно: %d", countIncrements)
}
