package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	cursorPrefix = "cursor"
)

// EncodeCursor uses base64 to encode a Node's field
func EncodeCursor(sortKey string, obj interface{}) string {
	rv := unwrapReflectValue(reflect.ValueOf(obj))
	for i := 0; i < rv.NumField(); i++ {
		structField := unwrapReflectType(rv.Type()).Field(i)
		field := unwrapReflectValue(reflect.Indirect(rv.Field(i)))
		fieldKind := field.Kind()
		fieldType := structField.Type
		fieldName := structField.Name
		if fieldKind == reflect.Struct && fieldType != TimeReflectType {
			res := EncodeCursor(sortKey, field.Interface())
			if len(res) > 0 {
				return res
			}
		}
		if strings.EqualFold(fieldName, sortKey) {
			switch fieldType {
			case reflect.TypeOf(""):
				stringCursor := field.String()
				return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s", cursorPrefix, stringCursor)))
			case TimeReflectType:
				timeCursor := field.Interface().(time.Time).Format("2006-01-02 15:04:05.000000")
				return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s", cursorPrefix, timeCursor)))
			}
		}
	}
	return ""
}

// DecodeCursor into raw value
func DecodeCursor(s string) (*string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New("invalid cursor provided")
	}
	prefixed := strings.TrimPrefix(string(decoded), cursorPrefix)
	return &prefixed, nil
}
