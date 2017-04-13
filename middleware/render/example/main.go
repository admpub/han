package main

import (
	"flag"

	"github.com/admpub/han"
	"github.com/admpub/han/engine/fasthttp"
	"github.com/admpub/han/engine/standard"
	mw "github.com/admpub/han/middleware"
	"github.com/admpub/han/middleware/render"
)

func main() {
	port := flag.String(`p`, "8080", "port")
	flag.Parse()
	e := han.New()
	e.Use(mw.Log())

	d := render.New(`standard`, `./template`)
	d.Init(true)

	e.Use(render.Middleware(d))

	e.Get("/", func(c han.Context) error {

		// It uses template file ./template/index.html
		return c.Render(`index`, map[string]interface{}{
			"Name": "Webx",
		})
	})

	// try visit: http://localhost:8080/api or http://localhost:8080/api?format=xml or
	// http://localhost:8080/api?format=json or
	// http://localhost:8080/api?format=jsonp&callback=f
	g := e.Group("/api", render.AutoOutput(nil))
	{
		g.Get("", func(c han.Context) error {
			c.Set("data", c.NewData().SetCode(1).SetData(han.H{
				"Name": "Webx",
			}))

			// It uses template file ./template/index.html
			c.Set("tmpl", "index")
			return nil
		})
	}

	switch `` {
	case `fast`:
		// FastHTTP
		e.Run(fasthttp.New(":" + *port))

	default:
		// Standard
		e.Run(standard.New(":" + *port))
	}
}
