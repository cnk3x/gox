package strs

import (
	"runtime"
	"strings"

	"github.com/cnk3x/gox/arrs"
)

var Windows = runtime.GOOS == "windows"

// Split splits a string into a slice of strings. return nil if the string is empty.
func Split(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

// Prefix returns the left part of the string before the first occurrence of the separator.
func Prefix(s string, sep string) string {
	if left, _, ok := strings.Cut(s, sep); ok {
		return left
	}
	return ""
}

// Prefix2 returns the left part of the string before the first occurrence of the separator.
func Prefix2(s1, s2 string, sep string) (string, string) { return Prefix(s1, sep), Prefix(s2, sep) }

// Suffix returns the right part of the string after the first occurrence of the separator.
func Suffix(s string, sep string) string { _, right, _ := strings.Cut(s, sep); return right }

// Suffix2 returns the right part of the string after the first occurrence of the separator.
func Suffix2(s1, s2 string, sep string) (string, string) { return Suffix(s1, sep), Suffix(s2, sep) }

// HasPrefix reports whether the string s begins with prefix. case-insensitive on Windows.
func HasPrefix(s string, prefix string) bool {
	return len(s) >= len(prefix) && Equal(s[:len(prefix)], prefix)
}

// TrimPrefix returns s without the provided leading prefix string. case-insensitive on Windows.
func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// HasSuffix reports whether the string s ends with suffix. case-insensitive on Windows.
func HasSuffix(s string, suffix string) bool {
	return len(s) >= len(suffix) && Equal(s[len(s)-len(suffix):], suffix)
}

// TrimSuffix returns s without the provided trailing suffix string. case-insensitive on Windows.
func TrimSuffix(s, suffix string) string {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// Equal returns true if a and b are equal, case-insensitive on Windows.
func Equal(a, b string) bool {
	if Windows {
		return strings.EqualFold(a, b)
	}
	return a == b
}

// PrefixMatch returns a function that returns true if the given string starts with any of the given prefixes.
func PrefixMatch(prefix ...string) func(string) bool {
	switch len(prefix) {
	case 0:
		return func(s string) bool { return false }
	case 1:
		return func(s string) bool { return HasPrefix(s, prefix[0]) }
	default:
		return func(s string) bool {
			return arrs.Some(prefix, func(it string) bool {
				return HasPrefix(s, it)
			})
		}
	}
}

// SuffixMatch returns a function that returns true if the given string ends with the given suffix.
func SuffixMatch(suffix string) func(string) bool {
	return func(s string) bool { return HasSuffix(s, suffix) }
}

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
// If n < 0, or len(n) == 0 there is no limit on the number of replacements.
func Replace(s, old, new string, n ...int) string {
	if len(n) == 0 {
		n = []int{-1}
	}
	return strings.Replace(s, old, new, n[0])
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
//
// Use Compare when you need to perform a three-way comparison (with
// [slices.SortFunc], for example). It is usually clearer and always faster
// to use the built-in string comparison operators ==, <, >, and so on.
func Compaare(a, b string) int { return strings.Compare(Lower(a), Lower(b)) }
