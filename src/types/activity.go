package types

import "time"

type ActivityListner interface {
	Listen(stop chan interface{}, events chan ActivityEvent) error
}

type ActivityEvent struct {
	LastActiveTimestamp time.Time
}
