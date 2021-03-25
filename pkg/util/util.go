package util

import (
	"context"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// TimeReflectType is exactly that
	TimeReflectType            = reflect.TypeOf(time.Time{})
)
// AddToContext adds a key/value pair on to context
func AddToContext(c *gin.Context, key ContextKey, value interface{}) *gin.Context {
	// return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
	ctx := context.WithValue(c.Request.Context(), key, value)
	c.Request = c.Request.WithContext(ctx)
	return c
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
