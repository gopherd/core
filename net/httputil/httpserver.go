package httputil

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/net/netutil"
)

var pong = []byte{'p', 'o', 'n', 'g'}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(pong)
}

type Config struct {
	Address           string            `json:"address"`
	StaticDir         string            `json:"static_dir"`
	StaticPath        string            `json:"static_path"`
	ConnTimeout       time.Duration     `json:"conn_timeout"`
	ReadHeaderTimeout time.Duration     `json:"read_header_timeout"`
	ReadTimeout       time.Duration     `json:"read_timeout"`
	WriteTimeout      time.Duration     `json:"write_timeout"`
	MaxConns          int               `json:"max_conns"`
	Headers           map[string]string `json:"headers"`
}

func (cfg *Config) autofix() {
	if cfg.ConnTimeout <= 0 {
		cfg.ConnTimeout = 60 * time.Second
	}
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 30 * time.Second
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 30 * time.Second
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 30 * time.Second
	}
	if cfg.MaxConns <= 0 {
		cfg.MaxConns = 1 << 15
	}
}

// HTTPServer ...
type HTTPServer struct {
	cfg    Config
	mux    http.ServeMux
	server *http.Server

	numHandling int64
}

func NewHTTPServer(cfg Config) *HTTPServer {
	cfg.autofix()
	httpd := &HTTPServer{
		cfg: cfg,
	}
	httpd.server = &http.Server{
		Addr:              httpd.cfg.Address,
		Handler:           &httpd.mux,
		ReadHeaderTimeout: httpd.cfg.ReadHeaderTimeout,
		ReadTimeout:       httpd.cfg.ReadTimeout,
		WriteTimeout:      httpd.cfg.WriteTimeout,
	}
	return httpd
}

func (httpd *HTTPServer) NumHandling() int64 {
	return atomic.LoadInt64(&httpd.numHandling)
}

func (httpd *HTTPServer) Addr() string {
	return httpd.server.Addr
}

func (httpd *HTTPServer) Listen() (net.Listener, error) {
	addr := httpd.server.Addr
	if addr == "" {
		addr = ":http"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if httpd.cfg.ConnTimeout > 0 {
		l = netutil.NewTCPKeepAliveListener(l.(*net.TCPListener), httpd.cfg.ConnTimeout)
	}
	if httpd.cfg.MaxConns <= 0 {
		return l, nil
	}
	return LimitListener(l, httpd.cfg.MaxConns), nil
}

func (httpd *HTTPServer) Serve(l net.Listener) error {
	return httpd.server.Serve(l)
}

func (httpd *HTTPServer) ListenAndServe() error {
	l, err := httpd.Listen()
	if err != nil {
		return err
	}
	return httpd.Serve(l)
}

func (httpd *HTTPServer) Shutdown(ctx context.Context) error {
	return httpd.server.Shutdown(ctx)
}

func (httpd *HTTPServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	httpd.Handle(pattern, http.HandlerFunc(handler), middlewares...)
}

func (httpd *HTTPServer) Handle(pattern string, handler http.Handler, middlewares ...Middleware) {
	for _, m := range middlewares {
		handler = m.Apply(handler)
	}
	httpd.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&httpd.numHandling, 1)
		defer atomic.AddInt64(&httpd.numHandling, -1)
		if httpd.cfg.Headers != nil {
			for k, v := range httpd.cfg.Headers {
				w.Header().Add(k, v)
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func (httpd *HTTPServer) JSONResponse(w http.ResponseWriter, r *http.Request, data any, options ...ResponseOptions) error {
	return JSONResponse(w, data, options...)
}
