package middlewares

import "net/http"

// Chain applies a series of middlewares to a given handler.
// It takes the handler and the middlewares as arguments.
// The middlewares are applied in reverse order, starting from the last middleware.
// The final handler with all the middlewares applied is returned.
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}
	return handler
}
