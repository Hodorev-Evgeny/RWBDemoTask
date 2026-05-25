package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
)

type SearchEvent struct {
	Query     string    `json:"query"`
	UserID    int64     `json:"user_id"`
	SessionID string    `json:"session_id"`
	TimeEvent time.Time `json:"time_event"`
}

func main() {
	var (
		natsURL = flag.String("url", "nats://localhost:4222", "NATS url")
		subject = flag.String("subject", "search.events", "NATS subject")
		total   = flag.Int("n", 100000, "total messages")
		workers = flag.Int("workers", 32, "workers count")
	)
	flag.Parse()

	nc, err := nats.Connect(*natsURL)
	if err != nil {
		log.Fatalf("connect nats: %v", err)
	}
	defer nc.Close()

	queries := []string{
		"iphone 15",
		"iphone 15",
		"iphone 15",
		"nike",
		"nike",
		"macbook",
		"playstation",
		"samsung",
	}

	jobs := make(chan int, *workers*4)

	var sent int64
	var failed int64

	start := time.Now()

	wg := sync.WaitGroup{}

	for w := 0; w < *workers; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range jobs {
				event := SearchEvent{
					Query:     queries[i%len(queries)],
					UserID:    int64(i + 1),
					SessionID: fmt.Sprintf("session-%d", i+1),
					TimeEvent: time.Now().UTC(),
				}

				body, err := json.Marshal(event)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					continue
				}

				if err := nc.Publish(*subject, body); err != nil {
					atomic.AddInt64(&failed, 1)
					continue
				}

				atomic.AddInt64(&sent, 1)
			}
		}()
	}

	for i := 0; i < *total; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()

	if err := nc.Flush(); err != nil {
		log.Fatalf("flush nats: %v", err)
	}

	duration := time.Since(start)

	fmt.Printf(
		"sent=%d failed=%d duration=%s messages_per_sec=%.2f\n",
		sent,
		failed,
		duration,
		float64(sent)/duration.Seconds(),
	)
}
