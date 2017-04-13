/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package render

import (
	"net/http"

	"github.com/admpub/han"
)

var (
	DefaultOptions = &Options{
		Skipper:              han.DefaultSkipper,
		DataKey:              `data`,
		TmplKey:              `tmpl`,
		DefaultTmpl:          `index`,
		JSONPCallbackName:    `callback`,
		OutputFunc:           Output,
		DefaultErrorHTTPCode: http.StatusInternalServerError,
	}
)

type Options struct {
	Skipper              han.Skipper
	DataKey              string
	TmplKey              string
	DefaultTmpl          string
	DefaultErrorTmpl     string
	JSONPCallbackName    string
	OutputFunc           func(format string, c han.Context, opt *Options) error
	DefaultErrorHTTPCode int
}

// Middleware set renderer
func Middleware(d han.Renderer) han.MiddlewareFunc {
	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			c.SetRenderer(d)
			return h.Handle(c)
		})
	}
}

func SetDefaultOptions(opt *Options) *Options {
	if opt.Skipper == nil {
		opt.Skipper = DefaultOptions.Skipper
	}
	if opt.OutputFunc == nil {
		opt.OutputFunc = DefaultOptions.OutputFunc
	}
	if len(opt.DataKey) == 0 {
		opt.DataKey = DefaultOptions.DataKey
	}
	if len(opt.TmplKey) == 0 {
		opt.TmplKey = DefaultOptions.TmplKey
	}
	if len(opt.DefaultTmpl) == 0 {
		opt.DefaultTmpl = DefaultOptions.DefaultTmpl
	}
	if len(opt.DefaultErrorTmpl) == 0 {
		opt.DefaultErrorTmpl = DefaultOptions.DefaultErrorTmpl
	}
	if len(opt.JSONPCallbackName) == 0 {
		opt.JSONPCallbackName = DefaultOptions.JSONPCallbackName
	}
	return opt
}

func checkOptions(options ...*Options) *Options {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	}
	if opt == nil {
		opt = DefaultOptions
	}
	return opt
}

// AutoOutput Outputs the specified format
func AutoOutput(options ...*Options) han.MiddlewareFunc {
	opt := checkOptions(options...)
	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			if opt.Skipper(c) {
				return h.Handle(c)
			}
			if err := h.Handle(c); err != nil {
				return err
			}
			return opt.OutputFunc(c.Format(), c, opt)
		})
	}
}

// Output Outputs the specified format
func Output(format string, c han.Context, opt *Options) error {
	switch format {
	case `json`:
		return c.JSON(c.Get(opt.DataKey))
	case `jsonp`:
		return c.JSONP(c.Query(opt.JSONPCallbackName), c.Get(opt.DataKey))
	case `xml`:
		return c.XML(c.Get(opt.DataKey))
	default:
		tmpl, ok := c.Get(opt.TmplKey).(string)
		if !ok {
			tmpl = opt.DefaultTmpl
		}
		data := c.Get(opt.DataKey)
		if v, y := data.(*han.Data); y {
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		if h, y := data.(han.H); y {
			v := h.ToData().SetContext(c)
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		return c.Render(tmpl, data)
	}
}

// SetFuncs register template function
func SetFuncs(c han.Context, v *han.Data) {
	c.SetFunc(`Info`, func() interface{} {
		return v.Info
	})
	c.SetFunc(`Code`, func() interface{} {
		return v.Code
	})
	c.SetFunc(`Zone`, func() interface{} {
		return v.Zone
	})
}

func HTTPErrorHandler(templates map[int]string, options ...*Options) han.HTTPErrorHandler {
	if templates == nil {
		templates = make(map[int]string)
	}
	tmplNum := len(templates)
	opt := checkOptions(options...)
	return func(err error, c han.Context) {
		code := opt.DefaultErrorHTTPCode
		var msg string
		if he, ok := err.(*han.HTTPError); ok {
			if he.Code > 0 {
				code = he.Code
			}
			msg = he.Message
		}
		title := http.StatusText(code)
		if c.Han().Debug() {
			msg = err.Error()
		} else if len(msg) == 0 {
			msg = title
		}
		if !c.Response().Committed() {
			switch {
			case c.Request().Method() == han.HEAD:
				c.NoContent(code)
			case tmplNum > 0:
				t, y := templates[code]
				if !y && code != 0 {
					t, y = templates[0]
				}
				if y {
					c.Set(opt.DataKey, c.NewData().SetInfo(han.H{
						"title":   title,
						"content": msg,
						"debug":   c.Han().Debug(),
						"code":    code,
					}))
					c.Set(opt.TmplKey, t)
					c.SetCode(code)
					if err := opt.OutputFunc(c.Format(), c, opt); err != nil {
						msg += "\n" + err.Error()
						y = false
						c.Han().Logger().Error(err)
					}
				}
				if y {
					break
				}
				fallthrough
			default:
				c.String(msg, code)
			}
		}
		c.Han().Logger().Debug(err)
	}
}
