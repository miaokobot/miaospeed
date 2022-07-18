package speed

import "sync"

type WriteCounter struct {
	Total     uint64
	RateLimit int64
	Lock      sync.Mutex
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Lock.Lock()
	wc.Total += uint64(n)
	wc.Lock.Unlock()
	return n, nil
}

func (wc *WriteCounter) Take() uint64 {
	wc.Lock.Lock()
	t := wc.Total
	wc.Total = 0
	wc.Lock.Unlock()
	return t
}
