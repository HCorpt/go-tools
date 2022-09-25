package hutils

import (
	"sync"
	"time"
)

type LocalTokenBucket struct {
	bufferTime   time.Duration
	timePerToken time.Duration
	cursor       time.Time
}

func NewLocalTokenBucket(rate int64, capacity int64) *LocalTokenBucket {
	if rate <= 0 || capacity <= 0 {
		return nil
	}
	return &LocalTokenBucket{
		bufferTime:   time.Duration(capacity * (int64(time.Second) / rate)),
		timePerToken: time.Second / time.Duration(rate),
		cursor:       time.Now(),
	}
}

// Take n tokens from Token Bucket
// if can't allocate token, return need time
// if success allocate token, return 0
func (tb *LocalTokenBucket) Take(n int64) time.Duration {
	now := time.Now()
	limitTime := now.Add(tb.bufferTime)
	if tb.cursor.Before(now) {
		tb.cursor = now
	}
	takedTime := tb.cursor.Add(time.Duration(n * int64(tb.timePerToken)))
	if takedTime.After(limitTime) {
		return takedTime.Sub(limitTime)
	}
	tb.cursor = takedTime
	return 0
}

type ThreadSafeTokenBucket struct {
	mtx sync.Mutex
	tb  *LocalTokenBucket
}

func NewThreadSafeTokenBucket(rate int64, capacity int64) *ThreadSafeTokenBucket {
	if rate <= 0 || capacity <= 0 {
		return nil
	}
	return &ThreadSafeTokenBucket{
		mtx: sync.Mutex{},
		tb:  NewLocalTokenBucket(rate, capacity),
	}
}

func (tb *ThreadSafeTokenBucket) Take(n int64) time.Duration {
	tb.mtx.Lock()
	defer tb.mtx.Unlock()
	return tb.tb.Take(n)
}
