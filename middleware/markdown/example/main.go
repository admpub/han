package main

import (
	"os"
	"path/filepath"

	"github.com/admpub/han"
	// "github.com/admpub/han/engine/fasthttp"
	"github.com/admpub/han/engine/standard"
	mw "github.com/admpub/han/middleware"
	"github.com/admpub/han/middleware/markdown"
)

func main() {
	e := han.New()
	e.Use(mw.Log(), mw.Recover())
	e.Use(markdown.Markdown(&markdown.Options{
		Path:   "/book/",
		Root:   filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/gopl-zh`),
		Browse: true,
	}))

	e.Get("/", han.HandlerFunc(func(c han.Context) error {
		return c.String("Hello, World!")
	}))

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
