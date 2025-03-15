package flags

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/cnk3x/gox/strs"
)

type StructOptions struct {
	EnvPrefix  string
	NamePrefix string
}

func EnvPrefix(prefix string) func(so *StructOptions) {
	return func(so *StructOptions) { so.EnvPrefix = prefix }
}

func NamePrefix(prefix string) func(so *StructOptions) {
	return func(so *StructOptions) { so.NamePrefix = prefix }
}

func FlagStruct(structObj any, options ...func(so *StructOptions)) FlagOption {
	var so StructOptions
	for _, fn := range options {
		fn(&so)
	}
	return func(fs *FlagSet) { addStruct(fs, structObj, so) }
}

func addStruct(fs *FlagSet, structObj any, so StructOptions) {
	rv := reflect.Indirect(reflect.ValueOf(structObj))
	rt := rv.Type()

	for i := range rt.NumField() {
		field := rt.Field(i)

		if !field.IsExported() {
			continue
		}

		fv, fk := rv.Field(i), field.Type.Kind()
		if fk == reflect.Pointer {
			if fv.IsNil() {
				fv.Set(reflect.New(field.Type.Elem()))
			}
			fv = fv.Elem()
			fk = field.Type.Elem().Kind()
		}

		if field.Anonymous {
			if fk == reflect.Struct {
				addStruct(fs, fv.Interface(), so)
			}
			continue
		}

		name := field.Tag.Get("flag")
		if name == "-" {
			continue
		}

		if name == "" {
			name = strs.Lower(field.Name)
		}

		if so.NamePrefix != "" {
			name = so.NamePrefix + "." + name
		}

		env := field.Tag.Get("env")
		if env == "" {
			env = strs.Replace(strs.Upper(name), "-", "_")
		}
		if env != "-" {
			if so.EnvPrefix != "" {
				env = so.EnvPrefix + "_" + env
			}
		}

		if fk == reflect.Struct {
			addStruct(fs, fv.Interface(), StructOptions{NamePrefix: name, EnvPrefix: env})
		} else {
			addFlag(fs, fv.Addr().Interface(), name, field.Tag.Get("short"), field.Tag.Get("usage"), env)
		}
	}
}

func addFlag(fs *FlagSet, val any, name, short, usage, env string) {
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
