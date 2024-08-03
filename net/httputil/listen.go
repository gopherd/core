package httputil

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gopherd/core/net/netutil"
)

// LimitListener returns a Listener that accepts at most n simultaneous
// connections from the provided Listener.
func LimitListener(l net.Listener, n int) net.Listener {
	return &limitListener{
		Listener: l,
		sem:      make(chan struct{}, n),
		done:     make(chan struct{}),
	}
}

type limitListener struct {
	net.Listener
	sem       chan struct{}
	closeOnce sync.Once     // ensures the done chan is only closed once
	done      chan struct{} // no values sent; closed when Close is called
}

// acquire acquires the limiting semaphore. Returns true if successfully
// accquired, false if the listener is closed and the semaphore is not
// acquired.
func (l *limitListener) acquire() bool {
	select {
	case <-l.done:
		return false
	case l.sem <- struct{}{}:
		return true
	}
}
func (l *limitListener) release() { <-l.sem }

func (l *limitListener) Accept() (net.Conn, error) {
	acquired := l.acquire()
	// If the semaphore isn't acquired because the listener was closed, expect
	// that this call to accept won't block, but immediately return an error.
	c, err := l.Listener.Accept()
	if err != nil {
		if acquired {
			l.release()
		}
		return nil, err
	}
	return &limitListenerConn{Conn: c, release: l.release}, nil
}

func (l *limitListener) Close() error {
	err := l.Listener.Close()
	l.closeOnce.Do(func() { close(l.done) })
	return err
}

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
}

// Listen creates a http server and registers a handler
//
// Wraps a netutil.ConnHandler handler as a websocket http.Handler as following:
//
//	import "golang.org/x/net/websocket"
//
//	websocket.Handler(func(conn *websocket.Conn) {
//		handler(netutil.IP(conn.Request()), conn)
//	})
func Listen(addr, path string, handler http.Handler, keepalive time.Duration) (*http.Server, net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}
	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	if keepalive > 0 {
		ln = netutil.NewTCPKeepAliveListener(ln.(*net.TCPListener), keepalive)
	}
	return server, ln, nil
}
