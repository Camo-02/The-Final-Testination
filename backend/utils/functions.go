package utils

func Map[T any, U any](vs []T, f func(T) U) []U {
	vsm := make([]U, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Filter[T any](vs []T, f func(T) bool) []T {
	vsf := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
