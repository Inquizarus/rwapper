package rwapper

import "net/http"

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}
	return handler
}
