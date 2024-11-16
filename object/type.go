package object

import "reflect"

func IsType(obj Object, name string) bool {
	type_ := reflect.TypeOf(obj)
	return type_.Kind() == reflect.Pointer && type_.Elem().Name() == name
}
