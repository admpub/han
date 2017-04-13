package middleware

import (
	"github.com/admpub/han"
)

// MaxAllowed limits simultaneous requests; can help with high traffic load
func MaxAllowed(n int) han.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			acquire() // before request
			err := h.Handle(c)
			release() // after request
			return err
		})
	}
}
