package strs

import (
	"cmp"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
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

func NumSet[T, R Ints | Uints | Floats](dst *R) (setter struct{ Set func(T, error) }) {
	setter.Set = func(val T, err error) {
		if err == nil {
			*dst = R(val)
		}
	}
	return
}

func SameSet[T any](dst *T) (setter struct{ Set func(T, error) }) {
	setter.Set = func(val T, err error) {
		if err == nil {
			*dst = val
		}
	}
	return
}

func ArrSet[T, R Ints | Uints | Floats](dst *[]R, valueAppend bool) (setter struct{ Set func(T, error) }) {
	setter.Set = func(t T, err error) {
		if err == nil {
			if valueAppend {
				*dst = append(*dst, R(t))
			} else {
				*dst = []R{R(t)}
			}
		}
	}
	return
}

func ArrSameSet[T any](dst *[]T, valueAppend bool) (setter struct{ Set func(T, error) }) {
	setter.Set = func(t T, err error) {
		if err == nil {
			if valueAppend {
				*dst = append(*dst, t)
			} else {
				*dst = []T{t}
			}
		}
	}
	return
}

func AnySet(dst any, val string, valueAppend bool) {
	switch x := dst.(type) {
	case *bool:
		SameSet(x).Set(strconv.ParseBool(val))
	case *string:
		*x = val
	case *int:
		SameSet(x).Set(strconv.Atoi(val))
	case *int8:
		NumSet[int64](x).Set(strconv.ParseInt(val, 0, 8))
	case *int16:
		NumSet[int64](x).Set(strconv.ParseInt(val, 0, 16))
	case *int32:
		NumSet[int64](x).Set(strconv.ParseInt(val, 0, 32))
	case *int64:
		NumSet[int64](x).Set(strconv.ParseInt(val, 0, 64))
	case *uint:
		NumSet[uint64](x).Set(strconv.ParseUint(val, 0, 64))
	case *uint8:
		NumSet[uint64](x).Set(strconv.ParseUint(val, 0, 8))
	case *uint16:
		NumSet[uint64](x).Set(strconv.ParseUint(val, 0, 16))
	case *uint32:
		NumSet[uint64](x).Set(strconv.ParseUint(val, 0, 32))
	case *uint64:
		NumSet[uint64](x).Set(strconv.ParseUint(val, 0, 64))
	case *float32:
		NumSet[float64](x).Set(strconv.ParseFloat(val, 32))
	case *float64:
		NumSet[float64](x).Set(strconv.ParseFloat(val, 32))
	case *net.IP:
		SameSet(x).Set(ParseIP(val))
	case *time.Duration:
		SameSet(x).Set(ParseDuration(val))
	case *[]bool:
		ArrSameSet(x, valueAppend).Set(strconv.ParseBool(val))
	case *[]string:
		if valueAppend {
			*x = append(*x, val)
		} else {
			*x = []string{val}
		}
	case *[]int:
		ArrSameSet(x, valueAppend).Set(strconv.Atoi(val))
	case *[]int32:
		ArrSet[int64](x, valueAppend).Set(strconv.ParseInt(val, 0, 32))
	case *[]int64:
		ArrSameSet(x, valueAppend).Set(strconv.ParseInt(val, 0, 64))
	case *[]uint:
		ArrSet[uint64](x, valueAppend).Set(strconv.ParseUint(val, 0, 32))
	case *[]float32:
		ArrSet[float64](x, valueAppend).Set(strconv.ParseFloat(val, 32))
	case *[]float64:
		ArrSameSet(x, valueAppend).Set(strconv.ParseFloat(val, 64))
	case *[]net.IP:
		ArrSameSet(x, valueAppend).Set(ParseIP(val))
	case *[]time.Duration:
		ArrSameSet(x, valueAppend).Set(ParseDuration(val))
	}
}

func ParseIP(val string) (net.IP, error) {
	if ip := net.ParseIP(val); ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("invalid ip address: %s", val)
}
