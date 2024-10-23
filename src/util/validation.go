package util

import "github.com/lib/pq"

// isUniqueViolation checks if the error is due to a unique constraint violation
func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // PostgreSQL unique violation code
	}
	return false
}
