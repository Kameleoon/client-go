package logging

import (
	"fmt"
	"reflect"
	"strings"
)

func ObjectToString(data interface{}) string {
	value := reflect.ValueOf(data)

	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf("'%s'", data)
	case reflect.Slice, reflect.Array:
		elements := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			elements = append(elements, ObjectToString(value.Index(i).Interface()))
		}
		return "[" + strings.Join(elements, ",") + "]"
	case reflect.Map:
		elements := make([]string, 0, value.Len())
		for _, key := range value.MapKeys() {
			keyString := ObjectToString(key.Interface())
			valString := ObjectToString(value.MapIndex(key).Interface())
			elements = append(elements, keyString+":"+valString)
		}
		return "{" + strings.Join(elements, ",") + "}"
	case reflect.Bool:
		if value.Bool() {
			return "true"
		}
		return "false"
	case reflect.Struct, reflect.Ptr:
		if value.Kind() == reflect.Ptr && value.IsNil() {
			return "<nil>"
		}
		if str, ok := data.(fmt.Stringer); ok {
			return str.String()
		}
	}
	return fmt.Sprintf("%#v", data)
}

func prepareArgs(args ...interface{}) []interface{} {
	for i := 0; i < len(args); i++ {
		args[i] = ObjectToString(args[i])
	}
	return args
}
