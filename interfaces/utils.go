package interfaces

func cloneSlice[T any](slice []T) []T {
	if slice == nil {
		return nil
	}

	return slice[:]
}
