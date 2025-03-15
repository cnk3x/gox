package strs

import (
	"cmp"
	"reflect"
	"strconv"
	"unsafe"
)

type (
	Strs interface{ ~string | ~[]byte }
	Ints interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}
	Uints interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}
	Floats interface{ ~float32 | ~float64 }
)

// Bytes 使用 unsafe 将 string 转换为 []byte
func Bytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func FromBytes(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func Int[T Ints](s string, def ...T) (r T) {
	i, err := strconv.ParseInt(s, 0, reflect.TypeOf(r).Bits())
	if err != nil {
		r = cmp.Or(def...)
	} else {
		r = T(i)
	}
	return
}

func FromInt[T Ints](i T, base int) string { return strconv.FormatInt(int64(i), base) }

func Uint[T Uints](s string, def ...T) (r T) {
	i, err := strconv.ParseUint(s, 0, reflect.TypeOf(r).Bits())
	if err != nil {
		r = cmp.Or(def...)
	} else {
		r = T(i)
	}
	return
}

func FromUint[T Uints](i T, base int) string { return strconv.FormatUint(uint64(i), base) }

func Float[T ~float32 | ~float64](s string, def ...T) (r T) {
	i, err := strconv.ParseFloat(s, reflect.TypeOf(r).Bits())
	if err != nil {
		r = cmp.Or(def...)
	} else {
		r = T(i)
	}
	return
}

func Bool[T ~bool](s string, def ...T) (r T) {
	i, err := strconv.ParseBool(s)
	if err != nil {
		r = cmp.Or(def...)
	} else {
		r = T(i)
	}
	return
}
