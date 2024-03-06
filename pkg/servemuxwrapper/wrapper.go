package servemuxwrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/inquizarus/rwapper/v2"
	"github.com/inquizarus/rwapper/v2/pkg/middlewares"
)

type wrapper struct {
	router *http.ServeMux
}

func (rw *wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

// Handle registers a handler for the specified method and path.
// It also applies any specified middlewares to the handler.
func (rw *wrapper) Handle(method string, path string, handler http.Handler, middlewaresList ...func(http.Handler) http.Handler) {
	v := strings.TrimLeft(fmt.Sprintf("%s %s", method, path), " ")
	rw.router.Handle(v, middlewares.Chain(handler, middlewaresList...))
}

func (rw *wrapper) Handler(method, path string, handler http.Handler, middlewaresList ...func(http.Handler) http.Handler) {
	rw.Handle(method, path, handler, middlewaresList...)
}

func (rw *wrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewaresList ...func(http.Handler) http.Handler) {
	rw.Handler(method, path, handler, middlewaresList...)
}

func (rw *wrapper) ParameterByName(name string, r *http.Request) string {
	return r.PathValue(strings.TrimRight(name, "."))
}

func (rw *wrapper) Parameterize(name string) string {
	return fmt.Sprintf("{%s}", strings.Trim(name, "{}"))
}

func New(router *http.ServeMux) rwapper.RouterWrapper {
	if router == nil {
		router = http.NewServeMux()
	}
	return &wrapper{
		router: router,
	}
}
