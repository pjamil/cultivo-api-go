package utils

import "time"

// TimePtr returns a pointer to the time.Time value.
func TimePtr(t time.Time) *time.Time {
	return &t
}
