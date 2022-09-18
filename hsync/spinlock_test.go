package hsync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func testLock(thread, times int, l sync.Locker) time.Duration {
	var wg sync.WaitGroup
	cnt1, cnt2 := 0, 0
	start := time.Now()
	for num := 0; num < thread; num += 1 {
		wg.Add(1)
		go func() {
			for n := 0; n < times; n += 1 {
				l.Lock()
				cnt1 += 1
				cnt2 += 2
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	end := time.Now()
	if cnt1 != thread*times {
		panic("count 1 with mismatch error")
	}
	if cnt2 != thread*times*2 {
		panic("count 2 with mismatch error")
	}
	return end.Sub(start)
}

func TestSpinLock(t *testing.T) {
	fmt.Printf("[1] spinlock %d ms\n", testLock(1, 1000000, &SpinLock{}).Milliseconds())
	fmt.Printf("[1] mutex    %d ms\n", testLock(1, 1000000, &sync.Mutex{}).Milliseconds())
	fmt.Printf("[2] spinlock %d ms\n", testLock(2, 1000000, &SpinLock{}).Milliseconds())
	fmt.Printf("[2] mutex    %d ms\n", testLock(2, 1000000, &sync.Mutex{}).Milliseconds())
	fmt.Printf("[4] spinlock %d ms\n", testLock(8, 1000000, &SpinLock{}).Milliseconds())
	fmt.Printf("[4] mutex    %d ms\n", testLock(8, 1000000, &sync.Mutex{}).Milliseconds())
	fmt.Printf("[64] spinlock %d ms\n", testLock(64, 1000000, &SpinLock{}).Milliseconds())
	fmt.Printf("[64] mutex    %d ms\n", testLock(64, 1000000, &sync.Mutex{}).Milliseconds())
}
