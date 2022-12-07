package utils

import (
	"reflect"
)

// GetFirstReturnType Receive a callback as an input and try to return the first return type as `reflect.Type` element.
// Return nil of the construct provided is not a function, is nil or do not return anything.
func GetFirstReturnType(construct any) reflect.Type {
	ctype := reflect.TypeOf(construct)

	if ctype == nil || ctype.Kind() != reflect.Func || ctype.NumOut() < 1 {
		return nil
	}

	return ctype.Out(0)
}

// EmbedsType checks that the provided `elem` interface embeds the provided type `e` directly. If it does, return true,
// otherwise return false.
func EmbedsType(elem interface{}, e reflect.Type) bool {
	if elem == nil {
		return false
	}

	etype, ok := elem.(reflect.Type)
	if !ok {
		etype = reflect.TypeOf(elem)
	}

	if etype.Kind() == reflect.Ptr {
		etype = etype.Elem()
	}

	for i := 0; i < etype.NumField(); i++ {
		field := etype.Field(i)

		if field.Anonymous && e == field.Type {
			return true
		}
	}

	return false
}
