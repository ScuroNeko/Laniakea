package utils

import (
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
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
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
		if fieldName == "-" {
			continue
		}

		// Handle omitempty
		isEmpty := field.IsZero()
		if slices.Contains(parts, "omitempty") && isEmpty {
			continue
		}

		var (
			fw  io.Writer
			err error
		)

		switch field.Kind() {
		case reflect.String:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(field.String()))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatInt(field.Int(), 10)))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatUint(field.Uint(), 10)))
			}
		case reflect.Float32:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatFloat(field.Float(), 'f', -1, 32)))
			}
		case reflect.Float64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatFloat(field.Float(), 'f', -1, 64)))
			}

		case reflect.Bool:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatBool(field.Bool())))
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Uint8 && !field.IsNil() {
				// Handle []byte as file upload (e.g., thumbnail)
				filename := fieldType.Tag.Get("filename")
				if filename == "" {
					filename = fieldName
				}
				fw, err = w.CreateFormFile(fieldName, filename)
				if err == nil {
					_, err = fw.Write(field.Bytes())
				}
			} else if !field.IsNil() {
				// Handle []string, []int, etc. — send as multiple fields with same name
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					fw, err = w.CreateFormField(fieldName)
					if err != nil {
						break
					}
					switch elem.Kind() {
					case reflect.String:
						_, err = fw.Write([]byte(elem.String()))
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						_, err = fw.Write([]byte(strconv.FormatInt(elem.Int(), 10)))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						_, err = fw.Write([]byte(strconv.FormatUint(elem.Uint(), 10)))
					case reflect.Bool:
						_, err = fw.Write([]byte(strconv.FormatBool(elem.Bool())))
					case reflect.Float32:
						_, err = fw.Write([]byte(strconv.FormatFloat(elem.Float(), 'f', -1, 32)))
					case reflect.Float64:
						_, err = fw.Write([]byte(strconv.FormatFloat(elem.Float(), 'f', -1, 64)))
					default:
						continue
					}
					if err != nil {
						break
					}
				}
			}

		case reflect.Struct:
			// Don't serialize structs as JSON — flatten them!
			// Telegram doesn't support nested JSON in form-data.
			// If you need nested data, use separate fields (e.g., ParseMode, CaptionEntities)
			// This is a design choice — you should avoid nested structs in params.
			return fmt.Errorf("nested structs are not supported in params — use flat fields")
		}

		if err != nil {
			return err
		}
	}

	return nil
}
