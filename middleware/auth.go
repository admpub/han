package middleware

import (
	"encoding/base64"
	"net/http"

	"github.com/admpub/han"
)

type (
	BasicValidateFunc func(string, string) bool
)

const (
	basic = "Basic"
)

// BasicAuth returns an HTTP basic authentication middleware.
//
// For valid credentials it calls the next handler.
// For invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicValidateFunc, skipper ...han.Skipper) han.MiddlewareFunc {
	var isSkiped han.Skipper
	if len(skipper) > 0 {
		isSkiped = skipper[0]
	} else {
		isSkiped = han.DefaultSkipper
	}
	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			if isSkiped(c) {
				return h.Handle(c)
			}
			auth := c.Request().Header().Get(han.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && auth[:l] == basic {
				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err == nil {
					cred := string(b)
					for i := 0; i < len(cred); i++ {
						if cred[i] == ':' {
							// Verify credentials
							if fn(cred[:i], cred[i+1:]) {
								return h.Handle(c)
							}
						}
					}
				}
			}
			c.Response().Header().Set(han.HeaderWWWAuthenticate, basic+" realm=Restricted")
			return han.NewHTTPError(http.StatusUnauthorized)
		})
	}
}
