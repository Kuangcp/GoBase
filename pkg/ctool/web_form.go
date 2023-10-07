package ctool

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	formTag        = "form"
	fmtTag         = "fmt"
	defaultTimeFmt = "2006-01-02"
)

// source: https://github.com/adonovan/gopl.io/blob/master/ch12/params/params.go

// Unpack populates the fields of the struct pointed to by ptr
// from the HTTP request parameters in req.
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	type Tags struct {
		val reflect.Value
		fmt string
	}
	// Build map of fields keyed by effective name.
	fieldMap := make(map[string]Tags)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get(formTag)
		// default use filed name with lower
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		_, ok := fieldMap[name]
		if ok {
			return fmt.Errorf("exist repeat field name: %s", name)
		}
		fieldMap[name] = Tags{val: v.Field(i), fmt: tag.Get(fmtTag)}
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fieldMap[name]
		if !f.val.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.val.Kind() == reflect.Slice {
				elem := reflect.New(f.val.Type().Elem()).Elem()
				if err := populate(elem, f.fmt, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.val.Set(reflect.Append(f.val, elem))
			} else {
				if err := populate(f.val, f.fmt, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, fmtStr, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Struct:
		switch v.Type().String() {
		case "time.Time":
			if err2 := handleTime(v, fmtStr, value, false); err2 != nil {
				return err2
			}
		default:
			return fmt.Errorf("unsupported struct %s", v.Type())
		}
		return nil
	case reflect.Pointer:
		switch v.Type().String() {
		case "*time.Time":
			if err2 := handleTime(v, fmtStr, value, true); err2 != nil {
				return err2
			}
		default:
			return fmt.Errorf("unsupported point %s", v.Type())
		}
		return nil
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

func handleTime(v reflect.Value, fmtStr, value string, point bool) error {
	if fmtStr == "" {
		fmtStr = defaultTimeFmt
	}
	parse, err := time.Parse(fmtStr, value)
	if err != nil {
		return err
	}
	if point {
		v.Set(reflect.ValueOf(&parse))
	} else {
		v.Set(reflect.ValueOf(parse))
	}
	return nil
}
