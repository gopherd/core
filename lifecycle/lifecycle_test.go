package lifecycle_test

import (
	"context"
	"testing"

	"github.com/gopherd/core/lifecycle"
)

func TestBaseLifecycle(t *testing.T) {
	var l lifecycle.BaseLifecycle
	if err := l.Init(context.Background()); err != nil {
		t.Errorf("BaseLifecycle.Init() error = %v, want nil", err)
	}
	if err := l.Start(context.Background()); err != nil {
		t.Errorf("BaseLifecycle.Start() error = %v, want nil", err)
	}
	if err := l.Shutdown(context.Background()); err != nil {
		t.Errorf("BaseLifecycle.Shutdown() error = %v, want nil", err)
	}
	if err := l.Uninit(context.Background()); err != nil {
		t.Errorf("BaseLifecycle.Uninit() error = %v, want nil", err)
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name string
		s    lifecycle.Status
		want string
	}{
		{
			name: "Created",
			s:    lifecycle.Created,
			want: "Created",
		},
		{
			name: "Starting",
			s:    lifecycle.Starting,
			want: "Starting",
		},
		{
			name: "Running",
			s:    lifecycle.Running,
			want: "Running",
		},
		{
			name: "Stopping",
			s:    lifecycle.Stopping,
			want: "Stopping",
		},
		{
			name: "Closed",
			s:    lifecycle.Closed,
			want: "Closed",
		},
		{
			name: "Unknown",
			s:    lifecycle.Status(100),
			want: "Unknown(100)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
