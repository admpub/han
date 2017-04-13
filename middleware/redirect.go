package middleware

import (
	"net/http"

	"github.com/admpub/han"
)

type (
	// RedirectConfig defines the config for Redirect middleware.
	RedirectConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper han.Skipper `json:"-"`

		// Status code to be used when redirecting the request.
		// Optional. Default value http.StatusMovedPermanently.
		Code int `json:"code"`
	}
)

var (
	// DefaultRedirectConfig is the default Redirect middleware config.
	DefaultRedirectConfig = RedirectConfig{
		Skipper: han.DefaultSkipper,
		Code:    http.StatusMovedPermanently,
	}
)

// HTTPSRedirect redirects HTTP requests to HTTPS.
// For example, http://webx.top will be redirect to https://webx.top.
//
// Usage `Han#Pre(HTTPSRedirect())`
func HTTPSRedirect() han.MiddlewareFuncd {
	return HTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func HTTPSRedirectWithConfig(config RedirectConfig) han.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			if !req.IsTLS() {
				host := req.Host()
				uri := req.URI()
				return c.Redirect("https://"+host+uri, config.Code)
			}
			return next.Handle(c)
		}
	}
}

// HTTPSWWWRedirect redirects HTTP requests to WWW HTTPS.
// For example, http://webx.top will be redirect to https://www.webx.top.
//
// Usage `Han#Pre(HTTPSWWWRedirect())`
func HTTPSWWWRedirect() han.MiddlewareFuncd {
	return HTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func HTTPSWWWRedirectWithConfig(config RedirectConfig) han.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			host := req.Host()
			uri := req.URI()
			if !req.IsTLS() && host[:4] != "www." {
				return c.Redirect("https://www."+host+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}

// WWWRedirect redirects non WWW requests to WWW.
// For example, http://webx.top will be redirect to http://www.webx.top.
//
// Usage `Han#Pre(WWWRedirect())`
func WWWRedirect() han.MiddlewareFuncd {
	return WWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func WWWRedirectWithConfig(config RedirectConfig) han.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			scheme := req.Scheme()
			host := req.Host()
			if host[:4] != "www." {
				uri := req.URI()
				return c.Redirect(scheme+"://www."+host+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}

// NonWWWRedirect redirects WWW requests to non WWW.
// For example, http://www.webx.top will be redirect to http://webx.top.
//
// Usage `Han#Pre(NonWWWRedirect())`
func NonWWWRedirect() han.MiddlewareFuncd {
	return NonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func NonWWWRedirectWithConfig(config RedirectConfig) han.MiddlewareFuncd {
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			scheme := req.Scheme()
			host := req.Host()
			if host[:4] == "www." {
				uri := req.URI()
				return c.Redirect(scheme+"://"+host[4:]+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}
