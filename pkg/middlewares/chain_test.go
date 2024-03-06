package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/rwapper/v3/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

const (
	testHeaderKey = "x-test-header"
	testValue     = "test"
)

func TestThatChainMiddlewaresWorksAsIntended(t *testing.T) {
	handler := middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testValue))
	}), func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(testHeaderKey, testValue)
			if nil != next {
				next.ServeHTTP(w, r)
			}
		})
	})

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "http://localhost", nil))

	assert.Equal(t, testValue, response.Body.String())
	assert.Equal(t, testValue, response.Header().Get(testHeaderKey))
}

func TestChainMiddlewareWithNoMiddlewares(t *testing.T) {
	// Test with no middlewares
	handler := middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testValue))
	}))

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "http://localhost", nil))

	assert.Equal(t, testValue, response.Body.String())
}

func TestChainMiddlewareWithMultipleMiddlewares(t *testing.T) {

	// Test with multiple middlewares
	handler := middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testValue))
	}),
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(testHeaderKey, testValue)
				if nil != next {
					next.ServeHTTP(w, r)
				}
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
				if nil != next {
					next.ServeHTTP(w, r)
				}
			})
		})

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "http://localhost", nil))

	assert.Equal(t, testValue, response.Body.String())
	assert.Equal(t, testValue, response.Header().Get(testHeaderKey))
	assert.Equal(t, http.StatusTeapot, response.Result().StatusCode)
}

func TestChainMiddlewareAppliesMiddlewaresInCorrectOrder(t *testing.T) {

	// Test the order of middleware execution
	var executedMiddleware string
	handler := middlewares.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testValue))
	}),
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				executedMiddleware += "First "
				if nil != next {
					next.ServeHTTP(w, r)
				}
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				executedMiddleware += "Second "
				if nil != next {
					next.ServeHTTP(w, r)
				}
			})
		})

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "http://localhost", nil))

	assert.Equal(t, "First Second ", executedMiddleware)
}
