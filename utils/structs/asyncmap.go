package structs

import "sync"

type AsyncMap[K Hashable, V any] struct {
	lock  sync.Mutex
	store map[K]V
}

func NewAsyncMap[K Hashable, T any]() *AsyncMap[K, T] {
	return &AsyncMap[K, T]{
		store: make(map[K]T),
	}
}

func (am *AsyncMap[K, T]) ForEach() map[K]T {
	am.lock.Lock()
	defer am.lock.Unlock()

	clones := make(map[K]T)
	for k, v := range am.store {
		clones[k] = v
	}

	return clones
}

func (am *AsyncMap[K, T]) Get(key K) (T, bool) {
	am.lock.Lock()
	defer am.lock.Unlock()

	val, ok := am.store[key]
	return val, ok
}

func (am *AsyncMap[K, T]) MustGet(key K) T {
	am.lock.Lock()
	defer am.lock.Unlock()

	return am.store[key]
}

func (am *AsyncMap[K, T]) Set(key K, val T) {
	am.lock.Lock()
	defer am.lock.Unlock()

	am.store[key] = val
}

func (am *AsyncMap[K, T]) Del(key K) {
	am.lock.Lock()
	defer am.lock.Unlock()

	delete(am.store, key)
}

func (am *AsyncMap[K, T]) Take(key K) (T, bool) {
	am.lock.Lock()
	defer am.lock.Unlock()

	val, ok := am.store[key]
	if ok {
		delete(am.store, key)
	}
	return val, ok
}
