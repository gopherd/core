package sms

import (
	"fmt"
	"sync"
)

type Provider interface {
	SendCode(phone, code string) error
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver func(source string) (Provider, error)

func Register(name string, driver Driver) {
	if driver == nil {
		panic("sms: Register driver is nil")
	}
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, dup := drivers[name]; dup {
		panic("sms: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(name, source string) (Provider, error) {
	driversMu.RLock()
	driver, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sms: unknown driver %q (forgotten import?)", name)
	}
	return driver(source)
}
