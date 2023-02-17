package rwapper

import "net/http"

type RouterWrapper interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Handle(method, path string, handler http.Handler, middlewares []Middleware)
	Handler(method, path string, handler http.Handler, middlewares []Middleware)
	HandlerFunc(method, path string, handler http.HandlerFunc, middlewares []Middleware)
	ParameterByName(name string, r *http.Request) string
	Parameterize(name string) string
}
