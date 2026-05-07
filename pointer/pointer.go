package pointer

// To returns the pointer to the given value.
//
//go:fix inline
func To[T any](value T) *T {
	return new(value)
}
