package flags

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/cnk3x/gox/strs"
)

type (
	StructOption  func(*StructOptions)
	StructOptions struct {
		EnvPrefix  string
		NamePrefix string
	}
)

func FlagStruct(structObj any, options ...StructOption) Option {
	return Flags(Struct(structObj, options...))
}

func EnvPrefix(prefix string) StructOption {
	return func(so *StructOptions) { so.EnvPrefix = prefix }
}

func NamePrefix(prefix string) StructOption {
	return func(so *StructOptions) { so.NamePrefix = prefix }
}

func Struct(structObj any, options ...StructOption) FlagOption {
	var so StructOptions
	for _, fn := range options {
		fn(&so)
	}
	return func(fs *FlagSet) { addStruct(fs, reflect.ValueOf(structObj), so) }
}

func addStruct(fs *FlagSet, rv reflect.Value, options StructOptions) {
	rv = reflect.Indirect(rv)
	rt := rv.Type()

	for i := range rt.NumField() {
		field := rt.Field(i)

		if !field.IsExported() || !typeSupport(field.Type, true) {
			continue
		}

		fv, ft := rv.Field(i), field.Type

		if !makeIfNil(rv.Field(i), ft) {
			continue
		}

		if field.Anonymous {
			if mayStruct(ft) {
				addStruct(fs, fv, options)
			}
			continue
		}

		if name, inline, short, usage, env, value := resloveTags(field, options); name != "-" {
			if mayStruct(ft) {
				if inline {
					addStruct(fs, fv, options)
				} else {
					addStruct(fs, fv, StructOptions{NamePrefix: name, EnvPrefix: env})
				}
			} else {
				if a, ok := getAddrInterface(fv); ok {
					if value != "" {
						strs.AnySet(a, value, false)
					}
					anyFlag(fs, a, name, short, usage, env)
				}
			}
		}
	}
}

func anyFlag(fs *FlagSet, val any, name, short, usage, env string) {
	if usage == "" {
		usage = strs.Replace(name, ".", " ")
	}

	if env != "-" && env != "" {
		usage += fmt.Sprintf(" (env: %s)", env)
		if envVal, envOk := os.LookupEnv(env); envOk {
			strs.AnySet(val, envVal, false)
		}
	}

	switch x := val.(type) {
	case *bool:
		fs.BoolVarP(x, name, short, *x, usage)
	case *string:
		fs.StringVarP(x, name, short, *x, usage)
	case *int:
		fs.IntVarP(x, name, short, *x, usage)
	case *int8:
		fs.Int8VarP(x, name, short, *x, usage)
	case *int16:
		fs.Int16VarP(x, name, short, *x, usage)
	case *int32:
		fs.Int32VarP(x, name, short, *x, usage)
	case *int64:
		fs.Int64VarP(x, name, short, *x, usage)
	case *uint:
		fs.UintVarP(x, name, short, *x, usage)
	case *uint8:
		fs.Uint8VarP(x, name, short, *x, usage)
	case *uint16:
		fs.Uint16VarP(x, name, short, *x, usage)
	case *uint32:
		fs.Uint32VarP(x, name, short, *x, usage)
	case *uint64:
		fs.Uint64VarP(x, name, short, *x, usage)
	case *float32:
		fs.Float32VarP(x, name, short, *x, usage)
	case *float64:
		fs.Float64VarP(x, name, short, *x, usage)
	case *net.IP:
		fs.IPVarP(x, name, short, *x, usage)
	case *time.Duration:
		fs.DurationVarP(x, name, short, *x, usage)
	case *[]bool:
		fs.BoolSliceVarP(x, name, short, *x, usage)
	case *[]string:
		fs.StringSliceVarP(x, name, short, *x, usage)
	case *[]int:
		fs.IntSliceVarP(x, name, short, *x, usage)
	case *[]int32:
		fs.Int32SliceVarP(x, name, short, *x, usage)
	case *[]int64:
		fs.Int64SliceVarP(x, name, short, *x, usage)
	case *[]uint:
		fs.UintSliceVarP(x, name, short, *x, usage)
	case *[]float32:
		fs.Float32SliceVarP(x, name, short, *x, usage)
	case *[]float64:
		fs.Float64SliceVarP(x, name, short, *x, usage)
	case *[]net.IP:
		fs.IPSliceVarP(x, name, short, *x, usage)
	case *[]time.Duration:
		fs.DurationSliceVarP(x, name, short, *x, usage)
	}
}

func resloveTags(field reflect.StructField, sOpt StructOptions) (name string, inline bool, short, usage, env, def string) {
	// name reslove
	if name = field.Tag.Get("flag"); name == "-" {
		return
	}

	if inline = name == ",inline"; inline {
		return
	}

	if name == "" {
		name = strs.Lower(field.Name)
	}

	// env reslove
	if env = field.Tag.Get("env"); env != "-" {
		if env == "" {
			env = strs.Replace(strs.Upper(name), "-", "_")
		}
		if sOpt.EnvPrefix != "" {
			env = sOpt.EnvPrefix + "_" + env
		}
	}

	// add name prefix
	if sOpt.NamePrefix != "" {
		name = sOpt.NamePrefix + "." + name
	}

	short, usage, def = field.Tag.Get("short"), field.Tag.Get("usage"), field.Tag.Get("default")
	return
}

func typeSupport(ft reflect.Type, checkElem bool) bool {
	switch ft.Kind() {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Struct:
		return true
	case reflect.Pointer:
		return checkElem && typeSupport(ft.Elem(), false)
	case reflect.Slice:
		return checkElem && typeSupport(ft.Elem(), false)
	default:
		return false
	}
}

func makeIfNil(v reflect.Value, t reflect.Type) bool {
	var isNil bool
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		isNil = v.IsNil()
	default:
		isNil = false
	}

	if isNil {
		if !v.CanSet() {
			return false
		}
		switch t.Kind() {
		case reflect.Pointer:
			v.Set(reflect.New(t.Elem()))
		case reflect.Slice:
			v.Set(reflect.MakeSlice(t, 0, 0))
		default:
			return false
		}
	}

	return true
}

func typeIndirect(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

func mayStruct(t reflect.Type) bool { return typeIndirect(t).Kind() == reflect.Struct }

func getAddrInterface(v reflect.Value) (any, bool) {
	if v.Kind() == reflect.Pointer {
		return v.Interface(), true
	}

	if v.CanAddr() {
		return v.Addr().Interface(), true
	}

	return nil, false
}
