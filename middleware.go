package rwapper

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}
	return handler
}
