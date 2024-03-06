package servemuxwrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/inquizarus/rwapper/v2"
)

type item struct {
	method      string
	handler     http.Handler
	middlewares []func(http.Handler) http.Handler
}

type container struct {
	items map[string][]item
}

func (ct *container) Add(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) {
	if items, ok := ct.items[path]; ok {
		items = append(items, item{
			method:      method,
			handler:     handler,
			middlewares: middlewares,
		})
		ct.items[path] = items
		return
	}

	ct.items[path] = []item{
		{
			method:      method,
			handler:     handler,
			middlewares: middlewares,
		},
	}

}

func (ct *container) PathIsRegistered(path string) bool {
	_, ok := ct.items[path]
	return ok
}

type wrapper struct {
	router *http.ServeMux
	ct     container
}

func (rw *wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

// Handle registers a handler for the specified method and path.
// It also applies any specified middlewares to the handler.
func (rw *wrapper) Handle(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) {
	// Check if the path is already registered
	if !rw.ct.PathIsRegistered(path) {
		// If the path is not registered, create a new handler function
		rw.router.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			// Check if the path has any registered items
			if items, ok := rw.ct.items[path]; ok {
				// Iterate over the registered items
				for _, item := range items {
					// Check if the method matches the current request method
					if item.method == r.Method {
						// Apply the registered middlewares and serve the request
						rwapper.ChainMiddleware(item.handler, item.middlewares...).ServeHTTP(w, r)
						return
					}
				}
			}
			// If no matching item is found, return a 404 status code
			w.WriteHeader(http.StatusNotFound)
		}))
	}

	// Add the method, path, handler, and middlewares to the router configuration
	rw.ct.Add(method, path, handler, middlewares...)
}

func (rw *wrapper) Handler(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) {
	rw.Handle(method, path, handler, middlewares...)
}

func (rw *wrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	rw.Handler(method, path, handler, middlewares...)
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
		ct: container{
			items: make(map[string][]item),
		},
	}
}
