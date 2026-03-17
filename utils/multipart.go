package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// Encode writes struct fields into multipart form-data using json tags as field names.
func Encode[T any](w *multipart.Writer, req T) error {
	v := unwrapMultipartValue(reflect.ValueOf(req))
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return fmt.Errorf("req must be a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldType.Name
		}

		parts := strings.Split(jsonTag, ",")
		fieldName := parts[0]
		if fieldName == "" {
			fieldName = fieldType.Name
		}
		if fieldName == "-" {
			continue
		}

		// Handle omitempty
		isEmpty := field.IsZero()
		if slices.Contains(parts, "omitempty") && isEmpty {
			continue
		}

		if err := writeMultipartField(w, fieldName, fieldType.Tag.Get("filename"), field); err != nil {
			return err
		}
	}

	return nil
}

func unwrapMultipartValue(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func writeMultipartField(w *multipart.Writer, fieldName, filename string, field reflect.Value) error {
	value := unwrapMultipartValue(field)
	if !value.IsValid() {
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		return writeMultipartValue(w, fieldName, []byte(value.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return writeMultipartValue(w, fieldName, []byte(strconv.FormatInt(value.Int(), 10)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return writeMultipartValue(w, fieldName, []byte(strconv.FormatUint(value.Uint(), 10)))
	case reflect.Float32:
		return writeMultipartValue(w, fieldName, []byte(strconv.FormatFloat(value.Float(), 'f', -1, 32)))
	case reflect.Float64:
		return writeMultipartValue(w, fieldName, []byte(strconv.FormatFloat(value.Float(), 'f', -1, 64)))
	case reflect.Bool:
		return writeMultipartValue(w, fieldName, []byte(strconv.FormatBool(value.Bool())))
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			if filename == "" {
				filename = fieldName
			}
			fw, err := w.CreateFormFile(fieldName, filename)
			if err != nil {
				return err
			}
			_, err = fw.Write(value.Bytes())
			return err
		}
	}

	// Telegram expects nested objects and arrays in multipart requests as JSON strings.
	data, err := json.Marshal(value.Interface())
	if err != nil {
		return err
	}
	if string(data) == "null" {
		return nil
	}
	return writeMultipartValue(w, fieldName, data)
}

func writeMultipartValue(w *multipart.Writer, fieldName string, value []byte) error {
	fw, err := w.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, strings.NewReader(string(value)))
	return err
}
