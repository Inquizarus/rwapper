package rwapper

import "net/http"

type RouterWrapper interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Handle(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler)
	Handler(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler)
	HandlerFunc(method, path string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler)
	ParameterByName(name string, r *http.Request) string
	Parameterize(name string) string
}
