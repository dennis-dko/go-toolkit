package datatype

import (
	"fmt"
	"strings"
)

func IsUtcTime(dateString string) bool {
	return strings.HasSuffix(dateString, "Z")
}

func SetAsUtc(dateString string) string {
	if IsUtcTime(dateString) {
		return dateString
	}
	return fmt.Sprintf("%sZ", dateString)
}
