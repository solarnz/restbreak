package types

import "time"

type Action struct {
	ActiveLimitDuration   time.Duration
	InactiveResetDuration time.Duration
	Command               string
	ResetOnLimit          bool
	RefireDuration        time.Duration
	RefireCount           int
	CurrentRefires        int

	LastInactive time.Time
	LastActive   time.Time
	LastFired    time.Time
}
