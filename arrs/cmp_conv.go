package arrs

func CompareConv[T, R any](compare func(a, b T) int, conv func(R) (T, bool)) func(a, b R) int {
	return func(a, b R) int {
		x, xok := conv(a)
		y, yok := conv(b)

		if !xok && !yok {
			return 0
		}

		if xok != yok {
			if xok {
				return -1
			}
			return 1
		}
		return compare(x, y)
	}
}

func LessConv[T, R any](less func(a, b T) bool, conv func(R) (T, bool)) func(a, b R) bool {
	return func(a, b R) bool {
		x, xok := conv(a)
		y, yok := conv(b)

		if !xok || !yok {
			return !xok
		}

		return less(x, y)
	}
}
