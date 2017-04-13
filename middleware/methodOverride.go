package middleware

import "github.com/admpub/han"

type (
	// MethodOverrideConfig defines the config for MethodOverride middleware.
	MethodOverrideConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper han.Skipper

		// Getter is a function that gets overridden method from the request.
		// Optional. Default values MethodFromHeader(han.HeaderXHTTPMethodOverride).
		Getter MethodOverrideGetter
	}

	// MethodOverrideGetter is a function that gets overridden method from the request
	MethodOverrideGetter func(han.Context) string
)

var (
	// DefaultMethodOverrideConfig is the default MethodOverride middleware config.
	DefaultMethodOverrideConfig = MethodOverrideConfig{
		Skipper: han.DefaultSkipper,
		Getter:  MethodFromHeader(han.HeaderXHTTPMethodOverride),
	}
)

// MethodOverride returns a MethodOverride middleware.
// MethodOverride  middleware checks for the overridden method from the request and
// uses it instead of the original method.
//
// For security reasons, only `POST` method can be overridden.
func MethodOverride() han.MiddlewareFuncd {
	return MethodOverrideWithConfig(DefaultMethodOverrideConfig)
}

// MethodOverrideWithConfig returns a MethodOverride middleware with config.
// See: `MethodOverride()`.
func MethodOverrideWithConfig(config MethodOverrideConfig) han.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultMethodOverrideConfig.Skipper
	}
	if config.Getter == nil {
		config.Getter = DefaultMethodOverrideConfig.Getter
	}

	return func(next han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			if req.Method() == han.POST {
				m := config.Getter(c)
				if m != "" {
					req.SetMethod(m)
				}
			}
			return next.Handle(c)
		}
	}
}

// MethodFromHeader is a `MethodOverrideGetter` that gets overridden method from
// the request header.
func MethodFromHeader(header string) MethodOverrideGetter {
	return func(c han.Context) string {
		return c.Request().Header().Get(header)
	}
}

// MethodFromForm is a `MethodOverrideGetter` that gets overridden method from the
// form parameter.
func MethodFromForm(param string) MethodOverrideGetter {
	return func(c han.Context) string {
		return c.Form(param)
	}
}

// MethodFromQuery is a `MethodOverrideGetter` that gets overridden method from
// the query parameter.
func MethodFromQuery(param string) MethodOverrideGetter {
	return func(c han.Context) string {
		return c.Query(param)
	}
}
