package helpers

import "reflect"

func IsEmpty(v interface{}) bool {
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !isZeroValue(field) {
			return false
		}
	}
	return true
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return IsEmpty(v.Interface())
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
