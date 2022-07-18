package memutils

import "time"

type MemDriver[T any] interface {
	Init(kargs ...string)

	Read(key string) (T, bool)
	Write(key string, value T, expire time.Duration, overwriteTTLIfExists bool) T
	IncBy(key string, value int, expire time.Duration, overwriteTTLIfExists bool) int
	Inc(key string, expire time.Duration, overwriteTTLIfExists bool) int

	List(prefix string) []string
	Expire(key string)
	SetExpire(key string, duration time.Duration) time.Duration
	Exists(key string) bool
	Wipe(prefix string)
	WipePrefix(prefix string)
}

func Now() int64 {
	return time.Now().UnixNano()
}

func Zero[T any]() T {
	var result T
	return result
}
