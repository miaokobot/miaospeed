package structs

type Set[T Hashable] struct {
	store map[T]bool
}

func (s *Set[T]) Has(key T) bool {
	_, ok := s.store[key]
	return ok
}

func (s *Set[T]) Add(key T) {
	s.store[key] = true
}

func (s *Set[T]) Remove(key T) {
	delete(s.store, key)
}

func (s *Set[T]) Digest() []T {
	arr := make([]T, len(s.store))
	i := 0
	for key := range s.store {
		arr[i] = key
		i++
	}
	return arr
}

func NewSet[T Hashable]() *Set[T] {
	return &Set[T]{
		store: make(map[T]bool),
	}
}
