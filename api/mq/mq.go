package mq

import (
	"fmt"
	"sync"

	"github.com/gopherd/core/api/discovery"
)

// Consumer used to consume received messages from mq.
type Consumer interface {
	// Setup would be called before consumption loop
	Setup() error
	// Cleanup would be called after consumption loop
	Cleanup() error
	// Consume runs the consumption loop blocking
	Consume(topic string, claim Claim)
}

// Claim used to receive requests from mq.
type Claim interface {
	// Err chan used to receive error
	Err() <-chan error
	// Message chan used to receive message content
	Message() <-chan []byte
}

// Conn is the top-level mq connection
type Conn interface {
	// Close closes the conn
	Close() error
	// Ping checks the connection to topic
	Ping(topic string) error
	// Subscribe subscribes topic with consumer
	Subscribe(topic string, consumer Consumer) error
	// Publish publishs message content to topic.
	// content allowed to be modified after Publish called.
	Publish(topic string, content []byte) error
}

// Driver is the interface that must be implemented by a mq driver
type Driver interface {
	// Open returns a Conn instance by a driver-specific source name
	Open(source string, discovery discovery.Discovery) (Conn, error)
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a mq driver available by the provided name
func Register(name string, driver Driver) {
	if driver == nil {
		panic("mq: Register driver is nil")
	}
	driversMu.Lock()
	defer driversMu.Unlock()
	if _, dup := drivers[name]; dup {
		panic("mq: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Open opens a mq connection specified by its driver name and
// a driver-specific source.
func Open(name, source string, discovery discovery.Discovery) (Conn, error) {
	driversMu.RLock()
	driver, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("mq: unknown driver %q (forgotten import?)", name)
	}
	return driver.Open(source, discovery)
}

// FuncConsumer implements Consumer interface
type FuncConsumer func(topic string, msg []byte, err error)

// Setup implements Consumer Setup method
func (fc FuncConsumer) Setup() error { return nil }

// Cleanup implements Consumer Cleanup method
func (fc FuncConsumer) Cleanup() error { return nil }

// Consume implements Consumer Consume method
func (fc FuncConsumer) Consume(topic string, claim Claim) {
	errChan := claim.Err()
	msgChan := claim.Message()
	for {
		select {
		case err := <-errChan:
			if err != nil {
				fc(topic, nil, err)
			}
			return
		case msg := <-msgChan:
			fc(topic, msg, nil)
		}
	}
}
