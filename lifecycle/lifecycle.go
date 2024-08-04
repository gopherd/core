package lifecycle

import "context"

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
