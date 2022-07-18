package memutils

import (
	"strings"
	"sync"
	"time"
)

type MemDriverMemory[T comparable] struct {
	MemDriver[T]

	mem   map[string]T
	timer map[string]int64

	lock sync.Mutex
}

func (md *MemDriverMemory[T]) Init(kargs ...string) {
	md.mem = make(map[string]T)
	md.timer = make(map[string]int64)
}

func (md *MemDriverMemory[T]) unsafeRead(key string) (T, bool) {
	now := Now()
	v, ok := md.mem[key]
	if ok {
		if t, ok := md.timer[key]; ok {
			if t <= now {
				delete(md.mem, key)
				delete(md.timer, key)
				return Zero[T](), false
			}
			return v, true
		} else {
			delete(md.mem, key)
			return Zero[T](), false
		}
	}
	return Zero[T](), false
}

func (md *MemDriverMemory[T]) Read(key string) (T, bool) {
	md.lock.Lock()
	defer md.lock.Unlock()

	return md.unsafeRead(key)
}

func (md *MemDriverMemory[T]) unsafeWrite(key string, value T, expire time.Duration, overwriteTTLIfExists bool) T {
	now := Now()
	_, ok := md.mem[key]
	md.mem[key] = value
	if !ok || overwriteTTLIfExists {
		md.timer[key] = now + expire.Nanoseconds()
	}

	return value
}

func (md *MemDriverMemory[T]) Write(key string, value T, expire time.Duration, overwriteTTLIfExists bool) T {
	md.lock.Lock()
	defer md.lock.Unlock()

	return md.unsafeWrite(key, value, expire, overwriteTTLIfExists)
}

func (md *MemDriverMemory[T]) IncBy(key string, value int, expire time.Duration, overwriteTTLIfExists bool) int {
	md.lock.Lock()
	defer md.lock.Unlock()

	val, ok := md.unsafeRead(key)
	nextVal := value

	if ok {
		switch any(val).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			nextVal += any(val).(int)
		default:
			return nextVal
		}
	}

	md.unsafeWrite(key, any(nextVal).(T), expire, overwriteTTLIfExists)
	return nextVal
}

func (md *MemDriverMemory[T]) Inc(key string, expire time.Duration, overwriteTTLIfExists bool) int {
	return md.IncBy(key, 1, expire, overwriteTTLIfExists)
}

func (md *MemDriverMemory[T]) Exists(key string) bool {
	md.lock.Lock()
	defer md.lock.Unlock()

	now := Now()
	if t, ok := md.timer[key]; ok && t > now {
		if _, ok := md.mem[key]; ok {
			return true
		}
	}
	return false
}

func (md *MemDriverMemory[T]) Expire(key string) {
	md.lock.Lock()
	defer md.lock.Unlock()

	_, ok := md.mem[key]
	if ok {
		delete(md.mem, key)
	}
	_, ok = md.timer[key]
	if ok {
		delete(md.timer, key)
	}
}

func (md *MemDriverMemory[T]) SetExpire(key string, duration time.Duration) time.Duration {
	md.lock.Lock()
	defer md.lock.Unlock()

	if _, ok := md.unsafeRead(key); ok {
		md.timer[key] = Now() + duration.Nanoseconds()
	}
	return duration
}

func (md *MemDriverMemory[T]) List(key string) []string {
	md.lock.Lock()
	defer md.lock.Unlock()

	slice := []string{}
	now := Now()
	for k, v := range md.timer {
		if now < v {
			slice = append(slice, k)
		}
	}
	return slice
}

func (md *MemDriverMemory[T]) Wipe(prefix string) {
	md.lock.Lock()
	defer md.lock.Unlock()

	md.Init()
}

func (md *MemDriverMemory[T]) WipePrefix(prefix string) {
	md.lock.Lock()
	defer md.lock.Unlock()

	for k := range md.mem {
		if strings.HasPrefix(k, prefix) {
			delete(md.mem, k)
			delete(md.timer, k)
		}
	}
}
