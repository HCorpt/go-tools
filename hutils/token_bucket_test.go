package hutils

import (
	"sync"
	"testing"
	"time"
)

func TestLocalTokenBucket(t *testing.T) {
	limiter := NewLocalTokenBucket(100, 10)
	end := time.Now().Add(time.Second * 30)
	unlimiteMap := make(map[int64]int)
	for time.Now().Before(end) {
		if limiter.Take(1) == 0 {
			roundUnix := time.Now().Round(time.Second).Unix()
			unlimiteMap[roundUnix] += 1
		}
		time.Sleep(time.Millisecond)
	}
	for _, cnt := range unlimiteMap {
		if cnt > 100 {
			t.Errorf("Token Bucket Unlimit Rate")
		}
	}
}

func TestThreadTokenBucket(t *testing.T) {
	limiter := NewThreadSafeTokenBucket(1000, 30)
	end := time.Now().Add(time.Second * 20)
	wg := &sync.WaitGroup{}
	result := make([]map[int64]int, 0)
	for n := 0; n < 10; n += 1 {
		wg.Add(1)
		result = append(result, make(map[int64]int))
		go func(m *map[int64]int, wg *sync.WaitGroup) {
			defer wg.Done()
			unlimiteMap := make(map[int64]int)
			for time.Now().Before(end) {
				if limiter.Take(1) == 0 {
					roundUnix := time.Now().Round(time.Second).Unix()
					unlimiteMap[roundUnix] += 1
				}
			}
		}(&result[n], wg)
	}
	wg.Wait()
	unlimiteMap := make(map[int64]int)
	for _, cnt := range result {
		for k, v := range cnt {
			unlimiteMap[k] += v
		}
	}
	for k, v := range unlimiteMap {
		if v > 1000 {
			t.Errorf("time %d break limit", k)
		}
	}
}
