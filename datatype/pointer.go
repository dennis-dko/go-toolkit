package datatype

import (
	"errors"
	"reflect"
)

func BoolPtr(b bool) *bool {
	return &b
}

func NullBoolPtr(nb NullBool) *NullBool {
	return &nb
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func NullFloat64Ptr(nf NullFloat64) *NullFloat64 {
	return &nf
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func NullInt64Ptr(ni NullInt64) *NullInt64 {
	return &ni
}

func IntPtr(i int) *int {
	return &i
}

func StringPtr(s string) *string {
	return &s
}

func NullStringPtr(ns NullString) *NullString {
	return &ns
}

func TimePtr(t CustomTime) *CustomTime {
	return &t
}

func NullTimePtr(nt NullTime) *NullTime {
	return &nt
}

func DatePtr(d CustomDate) *CustomDate {
	return &d
}

func NullDatePtr(nd NullDate) *NullDate {
	return &nd
}

// CheckPtrFieldValues checks if all fields of a struct are nil or zero value
// exceptFields is a list of fields that must be set
func CheckPtrFieldValues(rawStruct interface{}, exceptFields ...string) (*bool, error) {
	structData := reflect.ValueOf(rawStruct)
	if structData.Kind() == reflect.Ptr {
		structData = reflect.ValueOf(rawStruct).Elem()
	}
	if structData.Kind() != reflect.Struct {
		return nil, errors.New("invalid interface type or no data found")
	}
	fieldList := make(map[string]bool)
	for i := 0; i < structData.NumField(); i++ {
		if field := structData.Field(i); field.IsValid() {
			if (structData.Field(i).Kind() == reflect.Ptr && structData.Field(i).IsNil()) || structData.Field(i).IsZero() {
				fieldList[structData.Type().Field(i).Name] = true
			}
		}
	}
	if exceptFields != nil {
		exceptCount := 0
		for _, exceptField := range exceptFields {
			if _, ok := fieldList[exceptField]; !ok {
				exceptCount++
			}
		}
		return BoolPtr(exceptCount == len(exceptFields)), nil
	}
	return BoolPtr(len(fieldList) == structData.NumField()), nil
}
