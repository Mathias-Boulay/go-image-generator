package utils

/* Swap two values of the same type */
func Swap[T any](a, b T) (T, T) {
	return b, a
}
