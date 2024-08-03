package go_env

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var ErrIsNotStruct = errors.New("type is not struct")

func LoadFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		err := os.Setenv(key, value)
		if err != nil {
			return nil, err
		}
	}

	return Load[T]()
}

func Load[T any]() (*T, error) {
	var dest T

	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Struct {
		return nil, ErrIsNotStruct
	}

	v := reflect.ValueOf(&dest).Elem()
	numFields := t.NumField()

	for i := range numFields {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		tagValue := field.Tag.Get("env")
		if tagValue == "" {
			continue
		}

		items := strings.SplitN(tagValue, ",", 2)
		name := items[0]
		defaultValue := ""

		if len(items) == 2 {
			v := items[1]
			if strings.HasPrefix(v, "default:") {
				defaultValue = strings.TrimPrefix(v, "default:")
			}
		}

		value := os.Getenv(name)
		if value == "" {
			value = defaultValue
		}
		if value == "" {
			continue
		}

		kind := field.Type.Kind()
		log.Println(kind, value)

		switch kind {
		case reflect.Bool:
			r, err := strconv.ParseBool(value)
			if err != nil {
				return nil, parseErr(name, err)
			}
			v.FieldByIndex([]int{i}).SetBool(r)
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			bits := field.Type.Bits()

			r, err := strconv.ParseInt(value, 10, bits)
			if err != nil {
				return nil, parseErr(name, err)
			}

			v.FieldByIndex([]int{i}).SetInt(r)
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			r, err := strconv.ParseUint(value, 10, int(v.Type().Bits()))
			if err != nil {
				return nil, parseErr(name, err)
			}
			v.FieldByIndex([]int{i}).SetUint(r)
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			bits := field.Type.Bits()

			r, err := strconv.ParseFloat(value, bits)
			if err != nil {
				return nil, parseErr(name, err)
			}

			v.FieldByIndex([]int{i}).SetFloat(r)
		case reflect.String:
			v.FieldByIndex([]int{i}).SetString(value)
		case reflect.SliceOf(reflect.TypeFor[uint8]()).Kind():
			data, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return nil, err
			}

			v.FieldByIndex([]int{i}).SetBytes(data)
		default:
			return nil, fmt.Errorf("error unsupported type - %s", kind.String())
		}
	}

	return &dest, nil
}

func parseErr(name string, err error) error {
	return fmt.Errorf("error parsing %s - %w", name, err)
}
