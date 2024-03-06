package httprouterwrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/inquizarus/rwapper/v3"
	"github.com/inquizarus/rwapper/v3/pkg/middlewares"
	"github.com/julienschmidt/httprouter"
)

type httprouterWrapper struct {
	router *httprouter.Router
}

func (rw *httprouterWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

func (rw *httprouterWrapper) Handle(method, path string, handler http.Handler, middlewareList ...func(http.Handler) http.Handler) {
	rw.router.Handler(method, path, middlewares.Chain(handler, middlewareList...))
}

func (rw *httprouterWrapper) Handler(method, path string, handler http.Handler, middlewareList ...func(http.Handler) http.Handler) {
	rw.Handle(method, path, handler, middlewareList...)
}

func (rw *httprouterWrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewareList ...func(http.Handler) http.Handler) {
	rw.Handler(method, path, handler, middlewareList...)
}

func (rw *httprouterWrapper) ParameterByName(name string, r *http.Request) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(name)
}

func (rw *httprouterWrapper) Parameterize(name string) string {
	return fmt.Sprintf(":%s", strings.TrimLeft(name, ":"))
}

func New(router *httprouter.Router) rwapper.RouterWrapper {
	if router == nil {
		router = httprouter.New()
	}
	return &httprouterWrapper{
		router: router,
	}
}
