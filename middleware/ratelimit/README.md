## tollbooth_echo

[Han](https://github.com/admpub/han) middleware for rate limiting HTTP requests.


## Five Minutes Tutorial

```
package main

import (
	"time"

	"github.com/admpub/han"
	"github.com/admpub/han/engine/standard"
	"github.com/admpub/han/middleware/ratelimit"
)

func main() {
	e := han.New()

	// Create a limiter struct.
	limiter := ratelimit.NewLimiter(1, time.Second)

	e.Get("/", han.HandlerFunc(func(c han.Context) error {
		return c.String(200, "Hello, World!")
	}), ratelimit.LimitHandler(limiter))

	e.Run(standard.New(":4444"))
}

```