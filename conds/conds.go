package conds

func Iif[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}

func IifF[T any](c bool, t, f func() T) T {
	if c {
		return t()
	}
	return f()
}
