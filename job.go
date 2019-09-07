package main

import (
	"log"
	"sync"
	"time"
)

// TimeToLiveJob spawns a background job that expires URLs.
func TimeToLiveJob(cancel <-chan struct{}, wg *sync.WaitGroup, interval time.Duration, store Storage) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:
				store.DeleteSince(time.Now())
			case <-cancel:
				log.Println("Time-to-live Job shutting down ")
				return
			}
		}
	}()
}
