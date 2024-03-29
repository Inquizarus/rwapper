package chiwrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/inquizarus/rwapper/v3"
	"github.com/inquizarus/rwapper/v3/pkg/middlewares"
)

type chiRouterWrapper struct {
	router *chi.Mux
}

func (rw *chiRouterWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

func (rw *chiRouterWrapper) Handle(method, path string, handler http.Handler, middlewareList ...func(http.Handler) http.Handler) {
	rw.router.Method(method, path, middlewares.Chain(handler, middlewareList...))
}

func (rw *chiRouterWrapper) Handler(method, path string, handler http.Handler, middlewareList ...func(http.Handler) http.Handler) {
	rw.Handle(method, path, handler, middlewareList...)
}

func (rw *chiRouterWrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewareList ...func(http.Handler) http.Handler) {
	rw.Handler(method, path, handler, middlewareList...)
}

func (rw *chiRouterWrapper) ParameterByName(name string, r *http.Request) string {
	return chi.URLParam(r, name)
}

func (rw *chiRouterWrapper) Parameterize(name string) string {
	return fmt.Sprintf("{%s}", strings.Trim(name, "{}"))
}

func New(router *chi.Mux) rwapper.RouterWrapper {
	if router == nil {
		router = chi.NewRouter()
	}
	return &chiRouterWrapper{
		router: router,
	}
}
