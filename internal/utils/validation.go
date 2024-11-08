package utils

import (
	"errors"
	"github.com/lib/pq"
)

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // PostgreSQL unique violation code
	}
	return false
}
