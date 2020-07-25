package types

import (
	"time"

	"github.com/getlantern/systray"
)

type Action struct {
	Name string

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

	SystrayMenuItem *systray.MenuItem
}
