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

		formTags := strings.Split(fieldType.Tag.Get("json"), ",")
		fieldName := ""
		if len(formTags) == 0 {
			formTags = strings.Split(fieldType.Tag.Get("json"), ",")
		}

		if len(formTags) > 0 {
			fieldName = formTags[0]
			if fieldName == "-" {
				continue
			}
			if slices.Index(formTags, "omitempty") >= 0 {
				if field.IsZero() {
					continue
				}
			}
		} else {
			fieldName = strings.ToLower(fieldType.Name)
		}

		var (
			fw  io.Writer
			err error
		)
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				fw, err = w.CreateFormField(fieldName)
				if err == nil {
					_, err = fw.Write([]byte(field.String()))
				}
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
		case reflect.Float32, reflect.Float64:
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
				filename := fieldType.Tag.Get("filename")
				if filename == "" {
					filename = fieldName
				}

				ext := ""
				filename = filename + ext

				fw, err = w.CreateFormFile(fieldName, filename)
				if err == nil {
					_, err = fw.Write(field.Bytes())
				}
			} else if !field.IsNil() {
				// Handle slice of primitive values (as multiple form fields with the same name)
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					fw, err = w.CreateFormField(fieldName)
					if err == nil {
						switch elem.Kind() {
						case reflect.String:
							_, err = fw.Write([]byte(elem.String()))
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							_, err = fw.Write([]byte(strconv.FormatInt(elem.Int(), 10)))
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							_, err = fw.Write([]byte(strconv.FormatUint(elem.Uint(), 10)))
						}
					}
				}
			}

		case reflect.Struct:
			var jsonData []byte
			jsonData, err = json.Marshal(field.Interface())
			if err == nil {
				fw, err = w.CreateFormField(fieldName)
				if err == nil {
					_, err = fw.Write(jsonData)
				}
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
