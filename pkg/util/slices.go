package util

// Map can be used for map-reduce chaining on slices, since it's missing from the `slices` package
func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
