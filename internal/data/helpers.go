package data

import (
	"github.com/lib/pq"
	"github.com/google/uuid"
)

func isForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503"
	}
	return false
}

func GenerateUUID() string {
	return uuid.New().String()
}
