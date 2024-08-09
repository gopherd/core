package builtin

import (
	"context"
	"net/http"

	"github.com/gopherd/core/component"
)

// HTTPServerComponent is the interface that wraps the Handle method.
type HTTPServerComponent interface {
	Handle(path string, handler http.Handler)
}

const HTTPServerComponentName = "github.com/gopherd/core/component/httpserver"

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
		Addr            string
		DefaultServeMux bool
	}]
	mux    *http.ServeMux
	server *http.Server
}

func (com *httpServerComponent) Start(ctx context.Context) error {
	if com.Options().DefaultServeMux {
		com.mux = http.DefaultServeMux
	} else {
		com.mux = http.NewServeMux()
	}
	com.server = &http.Server{Addr: com.Options().Addr, Handler: com.mux}
	go com.server.ListenAndServe()
	return nil
}

func (com *httpServerComponent) Shutdown(ctx context.Context) error {
	return com.server.Shutdown(ctx)
}

// Handle implements the HTTPServerComponent interface.
func (com *httpServerComponent) Handle(path string, handler http.Handler) {
	com.mux.Handle(path, handler)
}
