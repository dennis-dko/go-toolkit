package util

import (
	"github.com/google/uuid"
)

// SetUUID generates a new UUID
func SetUUID() string {
	id := uuid.New()
	return id.String()
}
