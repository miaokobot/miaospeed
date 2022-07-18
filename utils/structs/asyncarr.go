package structs

import "sync"

type AsyncArr[T any] struct {
	lock  sync.Mutex
	store []T
}

func NewAsyncArr[T any]() *AsyncArr[T] {
	return &AsyncArr[T]{
		store: make([]T, 0),
	}
}

func (aa *AsyncArr[T]) ForEach() []T {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	return aa.store[:]
}

func (aa *AsyncArr[T]) Get(idx int) (*T, bool) {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	if idx >= 0 && idx < len(aa.store) {
		return &aa.store[idx], true
	}
	return nil, false
}

func (aa *AsyncArr[T]) MustGet(idx int) *T {
	t, _ := aa.Get(idx)
	return t
}

func (aa *AsyncArr[T]) Set(idx int, val T) bool {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	if idx >= 0 && idx < len(aa.store) {
		aa.store[idx] = val
		return true
	}
	return false
}

func (aa *AsyncArr[T]) Push(val T) {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	aa.store = append(aa.store, val)
}

func (aa *AsyncArr[T]) Del(idx int) *T {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	if idx >= 0 && idx < len(aa.store) {
		del := aa.store[idx]
		aa.store = append(aa.store[:idx], aa.store[idx+1:]...)
		return &del
	}

	return nil
}

func (aa *AsyncArr[T]) Take(idx int) *T {
	return aa.Del(idx)
}
