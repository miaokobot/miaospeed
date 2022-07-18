package structs

import "fmt"

var X = fmt.Sprintf

func Contains[T comparable](source []T, target T) bool {
	for _, s := range source {
		if s == target {
			return true
		}
	}
	return false
}

func MapContains[T any, U comparable](source []T, mapper func(T) U, target U) bool {
	for _, s := range source {
		if mapper(s) == target {
			return true
		}
	}
	return false
}

func Map[T, U any](source []T, mapper func(T) U) []U {
	result := make([]U, len(source))

	for i := range source {
		result[i] = mapper(source[i])
	}

	return result
}

func Filter[T any](source []T, mapper func(T) bool) []T {
	result := make([]T, 0)

	for i := range source {
		if mapper(source[i]) {
			result = append(result, source[i])
		}
	}

	return result
}

func FilterMap[K Hashable, T any](source map[K]T, mapper func(K, T) bool) map[K]T {
	result := make(map[K]T)

	for k := range source {
		if mapper(k, source[k]) {
			result[k] = source[k]
		}
	}

	return result
}

func Index[T any](source []T, mapper func(T) bool) int {
	for i := range source {
		if mapper(source[i]) {
			return i
		}
	}

	return -1
}

func Exist[T any](source []T, mapper func(T) bool) bool {
	for i := range source {
		if mapper(source[i]) {
			return true
		}
	}

	return false
}

func ExistMap[K Hashable, T any](source map[K]T, mapper func(K, T) bool) bool {
	for k := range source {
		if mapper(k, source[k]) {
			return true
		}
	}

	return false
}

func MapToArr[K Hashable, T any](source map[K]T) []T {
	result := make([]T, 0)

	for k := range source {
		result = append(result, source[k])
	}

	return result
}

func MapToArrMap[K Hashable, T, U any](source map[K]T, mapper func(K, T) U) []U {
	result := make([]U, 0)

	for k := range source {
		result = append(result, mapper(k, source[k]))
	}

	return result
}

func ArrToMap[K Hashable, T, U any](source []T, mapper func(T, int) (K, U)) map[K]U {
	result := make(map[K]U)

	for i := range source {
		k, v := mapper(source[i], i)
		result[k] = v
	}

	return result
}

func Uniq[T any, H Hashable](source []T, mapper func(T) H) []T {
	result := make([]T, 0)
	set := NewSet[H]()

	for i := range source {
		hashKey := mapper(source[i])
		if !set.Has(hashKey) {
			result = append(result, source[i])
			set.Add(hashKey)
		}
	}

	return result
}

func Concat[T any](sources ...[]T) []T {
	result := make([]T, 0)

	for _, items := range sources {
		result = append(result, items...)
	}

	return result
}

func SafeIndex[T any](arr []*T, idx int) *T {
	if idx >= len(arr) {
		return nil
	}
	return arr[idx]
}

func WithIn[T Integer](t, a, b T) T {
	if t < a {
		return a
	} else if b < t {
		return b
	}
	return t
}

func Max[T Integer](a ...T) T {
	if len(a) == 0 {
		return 0
	}

	max := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}

	return max
}

func Min[T Integer](a ...T) T {
	if len(a) == 0 {
		return 0
	}

	min := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] < min {
			min = a[i]
		}
	}

	return min
}
