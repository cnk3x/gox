package arrs

func Filter[S ~[]V, V any](s S, predicate func(V) bool) (r S) {
	for _, v := range s {
		if predicate(v) {
			r = append(r, v)
		}
	}
	return
}

func Filters[S ~[]V, V any](s S, predicate func(V, int) bool) (r S) {
	for i, v := range s {
		if predicate(v, i) {
			r = append(r, v)
		}
	}
	return
}

func IndexFilter[S ~[]V, V any](s S, predicate func(int) bool) (r S) {
	for i, v := range s {
		if predicate(i) {
			r = append(r, v)
		}
	}
	return
}

func Map[R, V any](s []V, predicate func(V) R) (r []R) {
	for _, v := range s {
		r = append(r, predicate(v))
	}
	return
}

func Maps[R, V any](s []V, predicate func(V, int) R) (r []R) {
	for i, v := range s {
		r = append(r, predicate(v, i))
	}
	return
}

func IndexMap[R, V any](s []V, predicate func(int) R) (r []R) {
	for i := range s {
		r = append(r, predicate(i))
	}
	return
}

func FilterMap[R, V any](s []V, predicate func(V) (R, bool)) (r []R) {
	for _, v := range s {
		if v, ok := predicate(v); ok {
			r = append(r, v)
		}
	}
	return
}

func FilterMaps[R, V any](s []V, predicate func(V, int) (R, bool)) (r []R) {
	for i, v := range s {
		if v, ok := predicate(v, i); ok {
			r = append(r, v)
		}
	}
	return
}

func IndexFilterMap[R, V any](s []V, predicate func(int) (R, bool)) (r []R) {
	for i := range s {
		if v, ok := predicate(i); ok {
			r = append(r, v)
		}
	}
	return
}
