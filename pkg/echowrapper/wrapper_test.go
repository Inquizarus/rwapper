package echowrapper_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/rwapper/v2/pkg/echowrapper"

	"github.com/stretchr/testify/assert"
)

const (
	testHeaderKey      = "x-test-header"
	testHeaderValue    = "test"
	testParameterKey   = "foo"
	testParameterValue = "bar"
	bodyString         = "test body"
)

func TestThatTheWrapperWorksAsIntended(t *testing.T) {
	wrapper := echowrapper.New(nil)

	testCases := []struct {
		method      string
		routerPath  string
		requestPath string
		body        []byte
		expected    string
		headerKey   string
		headerVal   string
	}{
		{
			method:      http.MethodGet,
			routerPath:  "/foo/" + wrapper.Parameterize(testParameterKey),
			requestPath: "/foo/" + testParameterValue,
			body:        nil,
			expected:    testParameterValue,
			headerKey:   testHeaderKey,
			headerVal:   testHeaderValue,
		},
		{
			method:      http.MethodPost,
			routerPath:  "/foo/" + wrapper.Parameterize(testParameterKey),
			requestPath: "/foo/" + testParameterValue,
			body:        []byte("request body"),
			expected:    testParameterValue,
			headerKey:   testHeaderKey,
			headerVal:   testHeaderValue,
		},
		{
			method:      http.MethodPut,
			routerPath:  "/foo/" + wrapper.Parameterize(testParameterKey),
			requestPath: "/foo/" + testParameterValue,
			body:        []byte("request body"),
			expected:    testParameterValue,
			headerKey:   testHeaderKey,
			headerVal:   testHeaderValue,
		},
		{
			method:      http.MethodDelete,
			routerPath:  "/foo/" + wrapper.Parameterize(testParameterKey),
			requestPath: "/foo/" + testParameterValue,
			body:        []byte("request body"),
			expected:    testParameterValue,
			headerKey:   testHeaderKey,
			headerVal:   testHeaderValue,
		},
		{
			method:      http.MethodPatch,
			routerPath:  "/foo/" + wrapper.Parameterize(testParameterKey),
			requestPath: "/foo/" + testParameterValue,
			body:        []byte("request body"),
			expected:    testParameterValue,
			headerKey:   testHeaderKey,
			headerVal:   testHeaderValue,
		},
	}

	for _, tc := range testCases {
		// Register the route for the test case
		wrapper.HandlerFunc(tc.method, tc.routerPath, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(wrapper.ParameterByName(testParameterKey, r)))
		}, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(testHeaderKey, testHeaderValue)
				if nil != next {
					next.ServeHTTP(w, r)
				}
			})
		})

		response := httptest.NewRecorder()
		request := httptest.NewRequest(tc.method, "http://localhost"+tc.requestPath, bytes.NewReader(tc.body))

		wrapper.ServeHTTP(response, request)

		assert.Equal(t, tc.expected, response.Body.String())
		assert.Equal(t, tc.headerVal, response.Header().Get(tc.headerKey))
	}
}
