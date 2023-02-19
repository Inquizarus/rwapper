package echowrapper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/rwapper/v2/pkg/echowrapper"

	"github.com/stretchr/testify/assert"
)

const (
	testHeaderKey = "x-test-header"
	testValue     = "test"
)

func TestThatTheWrapperWorksAsIntended(t *testing.T) {
	wrapper := echowrapper.New(nil)
	parameterKey := "bar"

	wrapper.HandlerFunc(http.MethodGet, "/foo/"+wrapper.Parameterize(parameterKey), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(wrapper.ParameterByName(parameterKey, r)))
	}, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(testHeaderKey, testValue)
			if nil != next {
				next.ServeHTTP(w, r)
			}
		})
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost/foo/"+testValue, nil)

	wrapper.ServeHTTP(response, request)

	assert.Equal(t, testValue, response.Body.String())
	assert.Equal(t, testValue, response.Header().Get(testHeaderKey))
}
