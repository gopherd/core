// Package enums provides a way to register enum descriptors and lookup them by name.
package enums

import "sync"

// Descriptor is a struct that describes an enum type.
type Descriptor struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Members     []MemberDescriptor `json:"members"`
}

// MemberDescriptor is a struct that describes an enum member.
type MemberDescriptor struct {
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

var (
	descriptorsMu sync.RWMutex
	descriptors   = make(map[string]*Descriptor)
)

// RegisterDescriptor registers an enum descriptor.
func RegisterDescriptor(descriptor *Descriptor) {
	descriptorsMu.Lock()
	defer descriptorsMu.Unlock()
	if _, dup := descriptors[descriptor.Name]; dup {
		panic("enums: RegisterDescriptor called twice for descriptor " + descriptor.Name)
	}
	descriptors[descriptor.Name] = descriptor
}

// LookupDescriptor looks up an enum descriptor by name.
func LookupDescriptor(name string) *Descriptor {
	descriptorsMu.RLock()
	defer descriptorsMu.RUnlock()
	descriptor := descriptors[name]
	return descriptor
}
