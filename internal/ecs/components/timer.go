package components

import "time"

type Timer struct {
	Remaining time.Duration
	StartedAt time.Time
}
