package iters

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func filterMap[R, V any, K Ordered](s func(yield func(K, V) bool), predicate func(V, K) (R, bool)) (r []R) {
	for k, v := range s {
		if i, ok := predicate(v, k); ok {
			r = append(r, i)
		}
	}
	return
}

func FilterMap[R, V any, K Ordered, F func(V) (R, bool) | func(V, K) (R, bool)](s func(yield func(K, V) bool), predicate F) (r []R) {
	if f, ok := any(predicate).(func(V) (R, bool)); ok {
		return filterMap(s, func(v V, k K) (R, bool) { return f(v) })
	}

	if f, ok := any(predicate).(func(V, K) (R, bool)); ok {
		return filterMap(s, func(v V, k K) (R, bool) { return f(v, k) })
	}

	return
}

func FilterMapIndex[R, V any, K Ordered](s func(yield func(K, V) bool), predicate func(K) (R, bool)) (r []R) {
	return filterMap(s, func(v V, k K) (R, bool) { return predicate(k) })
}

func Filter[V any, K Ordered, F func(V) bool | func(V, K) bool](s func(yield func(K, V) bool), predicate F) (r []V) {
	if f, ok := any(predicate).(func(V) bool); ok {
		return filterMap(s, func(v V, k K) (V, bool) { return v, f(v) })
	}

	if f, ok := any(predicate).(func(V, K) bool); ok {
		return filterMap(s, func(v V, k K) (V, bool) { return v, f(v, k) })
	}

	return
}

func FilterIndex[V any, K Ordered](s func(yield func(K, V) bool), predicate func(K) bool) (r []V) {
	return filterMap(s, func(v V, k K) (V, bool) { return v, predicate(k) })
}

func Map[R, V any, K Ordered, F func(V) R | func(V, K) R](s func(yield func(K, V) bool), predicate F) (r []R) {
	if f, ok := any(predicate).(func(V) R); ok {
		return filterMap(s, func(v V, k K) (R, bool) { return f(v), true })
	}

	if f, ok := any(predicate).(func(V, K) R); ok {
		return filterMap(s, func(v V, k K) (R, bool) { return f(v, k), true })
	}

	return
}

func MapIndex[R, V any, K Ordered](s func(yield func(K, V) bool), predicate func(K) R) (r []R) {
	return filterMap(s, func(v V, k K) (R, bool) { return predicate(k), true })
}
