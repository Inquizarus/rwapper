package echowrapper

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/inquizarus/rwapper"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	echoParamsContextKeyName = contextKey("echoxoxo")
)

type echoRouterWrapper struct {
	router *echo.Echo
}

func (rw *echoRouterWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw.router.ServeHTTP(w, r)
}

func (rw *echoRouterWrapper) Handle(method, path string, handler http.Handler, middlewares []rwapper.Middleware) {
	rw.router.Add(method, path, func(c echo.Context) error {

		params := map[string]string{}
		names := c.ParamNames()
		values := c.ParamValues()

		for i, name := range names {
			params[name] = values[i]
		}

		rwapper.ChainMiddleware(handler, middlewares...).ServeHTTP(c.Response(), c.Request().WithContext(context.WithValue(c.Request().Context(), echoParamsContextKeyName, params)))
		return nil
	})
}

func (rw *echoRouterWrapper) Handler(method, path string, handler http.Handler, middlewares []rwapper.Middleware) {
	rw.Handle(method, path, handler, middlewares)
}

func (rw *echoRouterWrapper) HandlerFunc(method, path string, handler http.HandlerFunc, middlewares []rwapper.Middleware) {
	rw.Handler(method, path, handler, middlewares)
}

func (rw *echoRouterWrapper) ParameterByName(name string, r *http.Request) string {
	ctxValue := r.Context().Value(echoParamsContextKeyName)
	if ctxValue != nil {
		params := ctxValue.(map[string]string)
		return params[name]
	}
	return ""
}

func (rw *echoRouterWrapper) Parameterize(name string) string {
	return fmt.Sprintf(":%s", strings.TrimLeft(name, ":"))
}

func New(router *echo.Echo) rwapper.RouterWrapper {
	if router == nil {
		router = echo.New()
	}
	return &echoRouterWrapper{
		router: router,
	}
}
