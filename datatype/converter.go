package datatype

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// StringToInt64Ptr converts a string to *int64
func StringToInt64Ptr(s string) (*int64, error) {
	// Attempt to convert the string to an int64
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to int64 (%w)", err)
	}
	return &i, nil
}

// StringToBoolPtr converts a string to *bool
func StringToBoolPtr(s string) (*bool, error) {
	// Attempt to convert the string to a bool
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to bool (%w)", err)
	}
	return &b, nil
}

// CustomTimeToFormat converts the *CustomTime to another specified format
func CustomTimeToFormat(customTime *CustomTime, format *string) (*CustomTime, error) {
	if customTime == nil || format == nil {
		return nil, errors.New("no custom time / specified format exists")
	}
	formattedTime := customTime.Time.Format(*format)
	convTime, err := ParseTime(formattedTime, customTime.NoUtc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert custom time to new format (%w)", err)
	}
	return convTime, nil
}

// CustomDateToFormat converts the *CustomDate to another specified format
func CustomDateToFormat(customDate *CustomDate, format *string) (*CustomDate, error) {
	if customDate == nil || format == nil {
		return nil, errors.New("no custom date / specified format exists")
	}
	formattedDate := customDate.Time.Format(*format)
	convDate, err := ParseDate(formattedDate, customDate.NoUtc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert custom date to new format (%w)", err)
	}
	return convDate, nil
}
