package util

import (
	"context"
	"html"
	"reflect"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// TimeReflectType is exactly that
	TimeReflectType            = reflect.TypeOf(time.Time{})
	regexAlphanumeric          = regexp.MustCompile("^[a-zA-Z0-9]+$")
	regexAlphanumericWithSpace = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	regexEmail                 = regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
	// regexURL enforces http/https
	regexURL = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)?`)
)

// AddToContext adds a key/value pair on to context
func AddToContext(c *gin.Context, key ContextKey, value interface{}) *gin.Context {
	// return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
	ctx := context.WithValue(c.Request.Context(), key, value)
	c.Request = c.Request.WithContext(ctx)
	return c
}

// IsURL checks whether a given string is a valid URL
func IsURL(s string) (valid bool) {
	return regexURL.MatchString(s)
}

// IsEmail checks whether a given string is a valid Email
// address
func IsEmail(s string) (valid bool) {
	return regexEmail.MatchString(s)
}

// IsAlphaNumeric checks whether a given string only
// contains alphanumeric characters(w/ space if specified)
func IsAlphaNumeric(s string, withSpace bool) (valid bool) {
	if len(s) == 0 {
		return true
	}
	if withSpace {
		return regexAlphanumericWithSpace.MatchString(s)
	}
	return regexAlphanumeric.MatchString(s)
}

// CheckLen is a helper func to check whether a string
// falls between a specific range of characters
func CheckLen(s string, min int, max int) (valid bool) {
	var (
		mi int = min
		ma int = max
	)
	if min > max {
		ma = min
		mi = max
	}
	if len(s) < mi || len(s) > ma {
		return false
	}
	return true
}

// EscapeStruct traverses a struct and escapes
// string or string slice
func EscapeStruct(obj interface{}) {
	value := unwrapReflectValue(reflect.ValueOf(obj))
	if value.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < value.NumField(); i++ {
		originalField := reflect.Indirect(value.Field(i))
		fieldValue := unwrapReflectValue(value.Field(i))
		if originalField.Kind() == reflect.Struct && !originalField.Type().ConvertibleTo(TimeReflectType) {
			EscapeStruct(fieldValue.Addr().Interface())
		}
		fieldValue = getActualFieldValue(fieldValue)

		isStr := fieldValue.Kind() == reflect.String
		isArr := fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array
		if !isStr && !isArr {
			continue
		}
		switch reflect.Indirect(fieldValue).Kind() {
		case reflect.String:
			str := fieldValue.String()
			fieldValue.SetString(html.EscapeString(str))
			continue
		case reflect.Slice:
			sli := reflect.MakeSlice(fieldValue.Type(), fieldValue.Len(), fieldValue.Cap())
			for j := 0; j < sli.Len(); j++ {
				innerValue := originalField.Index(j)
				switch reflect.Indirect(innerValue).Kind() {
				case reflect.String:
					sli.Index(j).Set(reflect.ValueOf(html.EscapeString(innerValue.String())))
					continue
				}
			}
			fieldValue.Set(sli)
		}
	}
}

// UnescapeStruct traverses a struct and unescapes
// string or string slice
func UnescapeStruct(obj interface{}) {
	value := unwrapReflectValue(reflect.ValueOf(obj))
	if value.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < value.NumField(); i++ {
		originalField := reflect.Indirect(value.Field(i))
		fieldValue := unwrapReflectValue(value.Field(i))
		if fieldValue.Kind() == reflect.Struct && !fieldValue.Type().ConvertibleTo(TimeReflectType) {
			UnescapeStruct(fieldValue.Addr().Interface())
		}
		fieldValue = getActualFieldValue(fieldValue)

		isStr := fieldValue.Kind() == reflect.String
		isArr := fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array
		if !isStr && !isArr {
			continue
		}
		switch reflect.Indirect(fieldValue).Kind() {
		case reflect.String:
			str := fieldValue.String()
			fieldValue.SetString(html.UnescapeString(str))
			continue
		case reflect.Slice:
			sli := reflect.MakeSlice(fieldValue.Type(), fieldValue.Len(), fieldValue.Cap())
			for j := 0; j < sli.Len(); j++ {
				innerValue := originalField.Index(j)
				switch reflect.Indirect(innerValue).Kind() {
				case reflect.String:
					sli.Index(j).Set(reflect.ValueOf(html.UnescapeString(innerValue.String())))
					continue
				}
			}
			fieldValue.Set(sli)
		}
	}
}

// Unwraps custom types to reveal its primitive data type
func getActualFieldValue(v reflect.Value) reflect.Value {
	unwrap := reflect.Indirect(v)
	original := reflect.Indirect(v)
	if original.Kind() == reflect.Struct && !original.Type().ConvertibleTo(TimeReflectType) {
		for i := 0; i < original.Type().NumField(); i++ {
			newFieldType := unwrapReflectType(original.Type().Field(i).Type)
			unwrap = reflect.New(newFieldType)
			if original.Type() != reflect.Indirect(unwrap).Type() {
				unwrap = getActualFieldValue(unwrap)
			}
			if unwrap.IsValid() {
				return unwrap
			}
		}
	}
	return unwrap
}

// Continually unwrap until we get the pointer's underlying value
func unwrapReflectValue(rv reflect.Value) reflect.Value {
	cpy := reflect.Indirect(rv)
	for cpy.Kind() == reflect.Ptr {
		cpy = cpy.Elem()
	}
	return cpy
}

// Continually unwrap until we get the pointer's underlying value
func unwrapReflectType(rt reflect.Type) reflect.Type {
	cpy := reflect.Indirect(reflect.New(rt)).Type()
	for cpy.Kind() == reflect.Ptr {
		cpy = cpy.Elem()
	}
	return cpy
}
