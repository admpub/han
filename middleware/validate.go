package middleware

import "github.com/admpub/han"

func Validate(generator func() han.Validator, skipper ...han.Skipper) han.MiddlewareFunc {
	var skip han.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = han.DefaultSkipper
	}
	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			if skip(c) {
				return h.Handle(c)
			}
			c.SetValidator(generator())
			return h.Handle(c)
		})
	}
}
