package arrs

import (
	"slices"
)

// FindOrZero search an element in a slice based on a predicate. It returns the element if found or a zero value otherwise.
func Find[T any](s []T, predicate func(item T) bool) (r T) {
	for i := range s {
		if predicate(s[i]) {
			r = s[i]
			break
		}
	}
	return
}

func FindOr[T comparable](s []T, predicate func(item T) bool, fallback ...T) (r T) {
	for i := range s {
		if predicate(s[i]) {
			return s[i]
		}
	}

	for _, it := range fallback {
		if it != r {
			return it
		}
	}

	return
}

// Contains reports whether v is present in s.
func Contains[T comparable](s []T, element T) bool {
	return slices.Contains(s, element)
}

// reports whether at least oneelement e of s satisfies f(e).
func Some[E any](s []E, f func(E) bool) bool {
	return slices.ContainsFunc(s, f)
}

// Index returns the index of the first occurrence of v in s, or -1 if not present.
func Index[T comparable](s []T, element T) int {
	return slices.Index(s, element)
}

// At returns the element at the given index.
// It returns the zero value if the index is out of range.
// if index is negative, it counts from the end of the slice.
func At[T any](s []T, index int) (r T) {
	if index < 0 {
		index = len(s) + index
	}
	if index >= 0 && index < len(s) {
		return s[index]
	}
	return // zero value
}

// ReplaceOrAppend replaces the first element in s that satisfies find with value.
func ReplaceOrAppend[S ~[]E, E any](s S, value E, find func(E) bool) S {
	i := IndexFunc(s, find)
	if i > -1 {
		s[i] = value
	} else {
		s = append(s, value)
	}
	return s
}

func Each[E any](s []E, eachFn func(E, int)) {
	for i, item := range s {
		eachFn(item, i)
	}
}

func Walk[E any](s []E, walkFn func(E, int) bool) {
	for i, item := range s {
		if !walkFn(item, i) {
			break
		}
	}
}

func CleanFunc[E any](e []E, eq func(E, E) bool) []E {
	var n int
	for i, item := range e {
		if n == 0 || !Some(e[:n], func(exist E) bool { return eq(exist, item) }) {
			e[n], n = e[i], n+1
		}
	}
	clear(e[n:])
	return e[:n]
}
