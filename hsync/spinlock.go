package hsync

import (
	"runtime"
	"sync/atomic"
	"time"
)

type SpinLock struct {
	lock int32
}

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&s.lock, 0, 1) {
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	if !atomic.CompareAndSwapInt32(&s.lock, 1, 0) {
		panic("SpinLock Fatal Error: Called a UnLocked Mutex")
	}
}

func (s *SpinLock) TryLock() bool {
	if atomic.CompareAndSwapInt32(&s.lock, 0, 1) {
		return true
	}
	return false
}

func (s *SpinLock) TryLockWithMaxTimes(times int) bool {
	if times <= 0 {
		return false
	}
	for i := 0; i < times; i++ {
		if atomic.CompareAndSwapInt32(&s.lock, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (s *SpinLock) LockWithTimeOut(duration time.Duration) bool {
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		if atomic.CompareAndSwapInt32(&s.lock, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}
