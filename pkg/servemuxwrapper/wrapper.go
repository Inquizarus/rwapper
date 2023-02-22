package servemuxwrapper

import (
	"net/http"

	"github.com/inquizarus/rwapper/v2"
)

type httprouterWrapper struct {
	router *http.ServeMux
}

func (rw *httprouterWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

func (rw *httprouterWrapper) Handle(_, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) {
	rw.router.Handle(path, rwapper.ChainMiddleware(handler, middlewares...))
}

func (rw *httprouterWrapper) Handler(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) {
	rw.Handle(method, path, handler, middlewares...)
}

func (rw *httprouterWrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	rw.Handler(method, path, handler, middlewares...)
}

func (rw *httprouterWrapper) ParameterByName(name string, _ *http.Request) string {
	return name
}

func (rw *httprouterWrapper) Parameterize(name string) string {
	return name
}

func New(router *http.ServeMux) rwapper.RouterWrapper {
	if router == nil {
		router = http.NewServeMux()
	}
	return &httprouterWrapper{
		router: router,
	}
}
