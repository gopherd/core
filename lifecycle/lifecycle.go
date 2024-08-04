package lifecycle

import (
	"context"
	"strconv"
)

// Status represents lifecycle state
type Status int

const (
	Created Status = iota
	Starting
	Running
	Stopping
	Closed
)

func (state Status) String() string {
	switch state {
	case Created:
		return "Created"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	case Closed:
		return "Closed"
	default:
		return "Unknown(" + strconv.Itoa(int(state)) + ")"
	}
}

type Lifecycle interface {
	Init(context.Context) error
	Uninit(context.Context) error
	Start(context.Context) error
	Shutdown(context.Context) error
}

type BaseLifecycle struct {
}

func (l *BaseLifecycle) Init(_ context.Context) error {
	return nil
}

func (l *BaseLifecycle) Uninit(_ context.Context) error {
	return nil
}

func (l *BaseLifecycle) Start(_ context.Context) error {
	return nil
}

func (l *BaseLifecycle) Shutdown(_ context.Context) error {
	return nil
}
