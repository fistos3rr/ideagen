package data

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
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
