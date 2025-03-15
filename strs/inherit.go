package strs

import "strings"

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func Cut(s, sep string) (prefix, suffix string, ok bool) { return strings.Cut(s, sep) }

// Fields splits the string s around each instance of one or more consecutive white space
// characters, as defined by [unicode.IsSpace], returning a slice of substrings of s or an
// empty slice if s contains only white space.
func Fields(s string) []string { return strings.Fields(s) }

// FieldsFunc splits the string s at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of s. If all code points in s satisfy f(c) or the
// string is empty, an empty slice is returned.
//
// FieldsFunc makes no guarantees about the order in which it calls f(c)
// and assumes that f always returns the same value for a given c.
func FieldsFunc(s string, f func(rune) bool) []string { return strings.FieldsFunc(s, f) }

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(s string) string { return strings.TrimSpace(s) }

// TrimFunc returns a slice of the string s with all leading
// and trailing Unicode code points c satisfying f(c) removed.
func TrimFunc(s string, f func(rune) bool) string { return strings.TrimFunc(s, f) }

// IndexFunc returns the index into s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func IndexFunc(s string, f func(rune) bool) int { return strings.IndexFunc(s, f) }

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s string, substr string) int { return strings.Index(s, substr) }

// Join concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join(strs []string, sep string) string { return strings.Join(strs, sep) }

// Lower returns s with all Unicode letters mapped to their lower case.
func Lower(s string) string { return strings.ToLower(s) }

// Upper returns s with all Unicode letters mapped to their upper case.
func Upper(s string) string { return strings.ToUpper(s) }
