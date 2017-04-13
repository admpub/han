package middleware

import "github.com/admpub/han"

func FuncMap(funcMap map[string]interface{}, skipper ...han.Skipper) han.MiddlewareFunc {
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

			c.SetFunc(`Lang`, c.Lang)
			c.SetFunc(`T`, c.T)
			c.SetFunc(`Stored`, c.Stored)
			c.SetFunc(`Cookie`, c.Cookie)
			c.SetFunc(`Session`, c.Session)
			c.SetFunc(`Query`, c.Query)
			c.SetFunc(`Form`, c.Form)
			c.SetFunc(`QueryValues`, c.QueryValues)
			c.SetFunc(`FormValues`, c.FormValues)
			c.SetFunc(`Param`, c.Param)
			c.SetFunc(`Atop`, c.Atop)
			req := c.Request()
			c.SetFunc(`URL`, req.URL)
			c.SetFunc(`Header`, req.Header)
			c.SetFunc(`Flash`, c.Flash)
			if c.Han().FuncMap != nil {
				for name, function := range c.Han().FuncMap {
					c.SetFunc(name, function)
				}
			}
			if funcMap != nil {
				for name, function := range funcMap {
					c.SetFunc(name, function)
				}
			}
			return h.Handle(c)
		})
	}
}
