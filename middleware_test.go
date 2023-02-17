package rwapper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/rwapper"
	"github.com/stretchr/testify/assert"
)

const (
	testHeaderKey = "x-test-header"
	testValue     = "test"
)

func TestThatChainMiddlewaresWorksAsIntended(t *testing.T) {
	handler := rwapper.ChainMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
