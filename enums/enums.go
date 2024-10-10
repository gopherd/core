// Package enums provides a way to register enum descriptors and lookup them by name.
package enums

import (
	"errors"
	"sync"
)

// Descriptor is a struct that describes an enum type.
type Descriptor struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Members     []MemberDescriptor `json:"members"`
}

// MemberDescriptor is a struct that describes an enum member.
type MemberDescriptor struct {
	Name        string `json:"name"`
	Value       any    `json:"value"`
	Description string `json:"description"`
}

// Registry is a struct that holds a map of enum descriptors.
type Registry struct {
	descriptorsMu sync.RWMutex
	descriptors   map[string]*Descriptor
}

// Register registers an enum descriptor.
func (r *Registry) Register(descriptor *Descriptor) error {
	r.descriptorsMu.Lock()
	defer r.descriptorsMu.Unlock()
	if r.descriptors == nil {
		r.descriptors = make(map[string]*Descriptor)
	}
	if _, dup := r.descriptors[descriptor.Name]; dup {
		return errors.New("enums: Register called twice for descriptor " + descriptor.Name)
	}
	r.descriptors[descriptor.Name] = descriptor
	return nil
}

// Lookup looks up an enum descriptor by name.
func (r *Registry) Lookup(name string) *Descriptor {
	r.descriptorsMu.RLock()
	defer r.descriptorsMu.RUnlock()
	if r.descriptors == nil {
		return nil
	}
	return r.descriptors[name]
}

// DefaultRegistry is the default [Registry].
var DefaultRegistry = &Registry{}
