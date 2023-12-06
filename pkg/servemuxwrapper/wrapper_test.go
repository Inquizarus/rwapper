package servemuxwrapper_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"

	"github.com/stretchr/testify/assert"
)

const (
	testHeaderKey   = "x-test-header"
	testHeaderValue = "test"
	bodyString      = "test body"
)

func TestThatTheWrapperWorksAsIntended(t *testing.T) {
	wrapper := servemuxwrapper.New(nil)

	testCases := []struct {
		method     string
		routerPath string
		body       []byte
		expected   string
		headerKey  string
		headerVal  string
	}{
		{
			method:     http.MethodGet,
			routerPath: "/foo",
			body:       nil,
			expected:   "foobar",
			headerKey:  testHeaderKey,
			headerVal:  testHeaderValue,
		},
		{
			method:     http.MethodPost,
			routerPath: "/foo",
			body:       []byte("request body"),
			expected:   "foobar",
			headerKey:  testHeaderKey,
			headerVal:  testHeaderValue,
		},
		{
			method:     http.MethodPut,
			routerPath: "/foo",
			body:       []byte("request body"),
			expected:   "foobar",
			headerKey:  testHeaderKey,
			headerVal:  testHeaderValue,
		},
		{
			method:     http.MethodDelete,
			routerPath: "/foo",
			body:       []byte("request body"),
			expected:   "foobar",
			headerKey:  testHeaderKey,
			headerVal:  testHeaderValue,
		},
		{
			method:     http.MethodPatch,
			routerPath: "/foo",
			body:       []byte("request body"),
			expected:   "foobar",
			headerKey:  testHeaderKey,
			headerVal:  testHeaderValue,
		},
	}

	for _, tc := range testCases {
		routeHandler := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tc.expected))
		}

		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(testHeaderKey, testHeaderValue)
				if next != nil {
					next.ServeHTTP(w, r)
				}
			})
		}

		wrapper.HandlerFunc(tc.method, tc.routerPath, routeHandler, middleware)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(tc.method, "http://localhost"+tc.routerPath, bytes.NewReader(tc.body))

		wrapper.ServeHTTP(response, request)

		assert.Equal(t, tc.expected, response.Body.String())
		assert.Equal(t, tc.headerVal, response.Header().Get(tc.headerKey))
	}
}
