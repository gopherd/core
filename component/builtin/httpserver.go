package builtin

import (
	"context"
	"net"
	"net/http"

	"github.com/gopherd/core/component"
)

// HTTPServerComponent is the interface that wraps the Handle method.
type HTTPServerComponent interface {
	Handle(path string, handler http.Handler)
}

// HTTPServerComponentName is the unique identifier for the HTTPServerComponent.
const HTTPServerComponentName = "go/httpserver"

func init() {
	// Register the HTTPServerComponent implementation.
	component.Register(HTTPServerComponentName, func() component.Component {
		return new(httpServerComponent)
	})
}

// Ensure httpServerComponent implements HTTPServerComponent interface.
var _ HTTPServerComponent = (*httpServerComponent)(nil)

type httpServerComponent struct {
	component.BaseComponent[struct {
		Addr   string // Addr is the address to listen on.
		Block  bool   // Block indicates whether the Start method should block.
		NewMux bool   // NewMux indicates whether to create a new ServeMux.
	}]
	mux    *http.ServeMux
	server *http.Server
}

// Start implements the component.Component interface.
func (com *httpServerComponent) Start(ctx context.Context) error {
	addr := com.Options().Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if com.Options().NewMux {
		com.mux = http.NewServeMux()
	} else {
		com.mux = http.DefaultServeMux
	}
	com.server = &http.Server{Addr: addr, Handler: com.mux}
	if com.Options().Block {
		com.Logger().Info("http server started", "addr", addr)
		return com.server.Serve(ln)
	}
	go func() {
		com.Logger().Info("http server started", "addr", addr)
		com.server.Serve(ln)
	}()
	return nil
}

// Shutdown implements the component.Component interface.
func (com *httpServerComponent) Shutdown(ctx context.Context) error {
	return com.server.Shutdown(ctx)
}

// Handle implements the HTTPServerComponent interface.
func (com *httpServerComponent) Handle(path string, handler http.Handler) {
	com.mux.Handle(path, handler)
}
