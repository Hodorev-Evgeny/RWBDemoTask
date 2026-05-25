package core_domain

import (
	core_redis "RWBDwmoTask/internal/core/repository/redis"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Bucket struct {
	Start time.Time
	Query string
	Count int64
}

type TopItem struct {
	Query string `json:"query"`
	Count int64  `json:"count"`
}

type Storage struct {
	mu         sync.RWMutex
	timeLife   time.Duration
	bucketSize time.Duration
	buckets    []Bucket
	total      map[string]int64
	cachedTop  []TopItem
	redis      *core_redis.RedisClient
}

func NewStorage(
	timeLife time.Duration,
	bucketSize time.Duration,
	redis *core_redis.RedisClient,
) *Storage {
	return &Storage{
		mu:         sync.RWMutex{},
		timeLife:   timeLife,
		bucketSize: bucketSize,
		buckets:    make([]Bucket, 0, 60),
		total:      make(map[string]int64),
		cachedTop:  make([]TopItem, 0),
		redis:      redis,
	}
}

func CleanBuckets(store *Storage) {
	actual := store.buckets[:0]
	for _, bucket := range store.buckets {
		if time.Since(bucket.Start) > store.timeLife {
			store.total[bucket.Query] -= bucket.Count

			if store.total[bucket.Query] <= 0 {
				delete(store.total, bucket.Query)
			}

			continue
		}
		actual = append(actual, bucket)
	}

	store.buckets = actual
}

func Rebuild(store *Storage) {
	list := make([]TopItem, 0, len(store.total))
	for query, count := range store.total {
		list = append(list, TopItem{
			Query: query,
			Count: count,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Count == list[j].Count {
			return list[i].Query < list[j].Query
		}
		return list[i].Count > list[j].Count
	})

	store.cachedTop = list
}

func (s *Storage) Run() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		CleanBuckets(s)
		Rebuild(s)
		s.mu.Unlock()
	}
}

func (s *Storage) Add(
	userID int64,
	sessionID string,
	query string,
	count int64) {
	ctx, close := context.WithTimeout(context.Background(), 3*time.Second)
	defer close()

	ans, err := s.redis.Protect(ctx, userID, sessionID, query)
	if err != nil {
		fmt.Println("redis protect error:", err)
	}
	if !ans {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	CleanBuckets(s)

	s.buckets = append(s.buckets, Bucket{
		Start: time.Now(),
		Query: query,
		Count: count,
	})

	s.total[query] += count
	Rebuild(s)
}

func (s *Storage) GetTop(limit int) []TopItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		return []TopItem{}
	}

	if limit > len(s.cachedTop) {
		limit = len(s.cachedTop)
	}

	ans := make([]TopItem, limit)
	copy(ans, s.cachedTop[:limit])

	return ans
}
