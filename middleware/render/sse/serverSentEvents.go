package sse

import (
	"io"

	"github.com/admpub/sse"
	"github.com/admpub/han"
	"github.com/admpub/han/middleware/render"
	"github.com/admpub/han/middleware/render/driver"
)

func init() {
	render.Reg(`sse`, func(_ string) driver.Driver {
		return New()
	})
}

func New() *ServerSentEvents {
	return &ServerSentEvents{
		NopRenderer: &driver.NopRenderer{},
	}
}

type ServerSentEvents struct {
	*driver.NopRenderer
}

func (s *ServerSentEvents) Render(w io.Writer, name string, data interface{}, c han.Context) error {
	if v, y := data.(sse.Event); y {
		return sse.Encode(w, v)
	}
	return sse.Encode(w, sse.Event{
		Event: name,
		Data:  data,
	})
}
