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
	limiter := ratelimit.New(1, time.Second)

	e.Get("/", han.HandlerFunc(func(c han.Context) error {
		return c.String("Hello, World!")
	}), ratelimit.LimitHandler(limiter))

	e.Run(standard.New(":4444"))
}
