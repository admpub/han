package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/admpub/han"
)

type (
	// CORSConfig defines the config for CORS middleware.
	CORSConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper han.Skipper

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional with default value as []string{"*"}.
		AllowOrigins []string

		// AllowMethods defines a list methods allowed when accessing the resource.
		// This is used in response to a preflight request.
		// Optional with default value as `DefaultCORSConfig.AllowMethods`.
		AllowMethods []string

		// AllowHeaders defines a list of request headers that can be used when
		// making the actual request. This in response to a preflight request.
		// Optional with default value as []string{}.
		AllowHeaders []string

		// AllowCredentials indicates whether or not the response to the request
		// can be exposed when the credentials flag is true. When used as part of
		// a response to a preflight request, this indicates whether or not the
		// actual request can be made using credentials.
		// Optional with default value as false.
		AllowCredentials bool

		// ExposeHeaders defines a whitelist headers that clients are allowed to
		// access.
		// Optional with default value as []string{}.
		ExposeHeaders []string

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached.
		// Optional with default value as 0.
		MaxAge int
	}
)

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:      han.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{han.GET, han.HEAD, han.PUT, han.POST, han.DELETE},
	}
)

// CORS returns a cross-origin HTTP request (CORS) middleware.
// See https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func CORS() han.MiddlewareFunc {
	return CORSWithConfig(DefaultCORSConfig)
}

// CORSFromConfig returns a CORS middleware from config.
// See `CORS()`.
func CORSWithConfig(config CORSConfig) han.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCORSConfig.Skipper
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}
	allowOrigins := strings.Join(config.AllowOrigins, ",")
	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}
			req := c.Request()
			header := c.Response().Header()

			// Simple request
			if req.Method() != han.OPTIONS {
				header.Add(han.HeaderVary, han.HeaderOrigin)
				header.Set(han.HeaderAccessControlAllowOrigin, allowOrigins)
				if config.AllowCredentials {
					header.Set(han.HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					header.Set(han.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				return next.Handle(c)
			}

			// Preflight request
			header.Add(han.HeaderVary, han.HeaderOrigin)
			header.Add(han.HeaderVary, han.HeaderAccessControlRequestMethod)
			header.Add(han.HeaderVary, han.HeaderAccessControlRequestHeaders)
			header.Set(han.HeaderAccessControlAllowOrigin, allowOrigins)
			header.Set(han.HeaderAccessControlAllowMethods, allowMethods)
			if config.AllowCredentials {
				header.Set(han.HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				header.Set(han.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := req.Header().Get(han.HeaderAccessControlRequestHeaders)
				if h != "" {
					header.Set(han.HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge > 0 {
				header.Set(han.HeaderAccessControlMaxAge, maxAge)
			}
			return c.NoContent(http.StatusNoContent)
		})
	}
}
