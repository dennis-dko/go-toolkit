package util

import (
	"fmt"
	"reflect"
)

func StringifyMap(data interface{}) string {
	if reflect.TypeOf(data).Kind() != reflect.Map {
		return ""
	}
	return fmt.Sprint(data)
}
