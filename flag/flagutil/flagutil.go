package flagutil

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

type FlagSet interface {
	// define field as a flag:
	//
	// N int `flag:"n", help:"help information"`
	//
	// supported scalar types:
	//	int,int32,int64,uint,uint32,uint64,float32,float64,bool,string,time.Duration
}

func Parse(name string, output io.Writer, fset FlagSet, args []string) ([]string, error) {
	s := flag.NewFlagSet(name, flag.ContinueOnError)
	if output == nil {
		s.SetOutput(io.Discard)
	} else {
		s.SetOutput(output)
	}
	if err := addFlags(s, reflect.StructField{}, reflect.ValueOf(fset)); err != nil {
		return nil, err
	}
	if err := s.Parse(args); err != nil {
		return nil, err
	}
	return s.Args(), nil
}

func ParseFlags(fset FlagSet, args []string) ([]string, error) {
	return Parse("", nil, fset, args)
}

// addFlags scans fields of structs recursively to find things with flag tags
// and add them to the flag set.
func addFlags(f *flag.FlagSet, field reflect.StructField, value reflect.Value) error {
	// is it a field we are allowed to reflect on?
	if field.PkgPath != "" {
		return nil
	}
	// now see if is actually a flag
	flagName, isFlag := field.Tag.Lookup("flag")
	help := field.Tag.Get("help")
	if !isFlag {
		// not a flag, but it might be a struct with flags in it
		if value.Elem().Kind() != reflect.Struct {
			return nil
		}
		// go through all the fields of the struct
		sv := value.Elem()
		for i := 0; i < sv.Type().NumField(); i++ {
			child := sv.Type().Field(i)
			v := sv.Field(i)
			// make sure we have a pointer
			if v.Kind() != reflect.Ptr {
				v = v.Addr()
			}
			// check if that field is a flag or contains flags
			if err := addFlags(f, child, v); err != nil {
				return err
			}
		}
		return nil
	}
	switch v := value.Interface().(type) {
	case flag.Value:
		f.Var(v, flagName, help)
	case *bool:
		f.BoolVar(v, flagName, *v, help)
	case *time.Duration:
		f.DurationVar(v, flagName, *v, help)
	case *float32:
		f.Var(NewFloat32(*v, v), flagName, help)
	case *float64:
		f.Float64Var(v, flagName, *v, help)
	case *int:
		f.IntVar(v, flagName, *v, help)
	case *int32:
		f.Var(NewInt32(*v, v), flagName, help)
	case *int64:
		f.Int64Var(v, flagName, *v, help)
	case *string:
		f.StringVar(v, flagName, *v, help)
	case *uint:
		f.UintVar(v, flagName, *v, help)
	case *uint32:
		f.Var(NewUint32(*v, v), flagName, help)
	case *uint64:
		f.Uint64Var(v, flagName, *v, help)
	default:
		return fmt.Errorf("Cannot understand flag of type %T", v)
	}
	return nil
}

// errParse is returned by Set if a flag's value fails to parse, such as with an invalid integer for Int.
// It then gets wrapped through failf to provide more information.
var errParse = errors.New("parse error")

// errRange is returned by Set if a flag's value is out of range.
// It then gets wrapped through failf to provide more information.
var errRange = errors.New("value out of range")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}
	return err
}
