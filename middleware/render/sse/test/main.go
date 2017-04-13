package main

import (
	"fmt"
	"math/rand"

	"github.com/admpub/han"
	"github.com/admpub/han/engine/fasthttp"
	"github.com/admpub/han/engine/standard"
	mw "github.com/admpub/han/middleware"
	"github.com/admpub/han/middleware/render"
	_ "github.com/admpub/han/middleware/render/sse"
)

func main() {
	engine := ``
	e := han.New()
	e.Use(mw.Log(), mw.Recover())
	e.Use(render.Middleware(render.New(`sse`, ``)))

	e.Get("/room/:roomid", roomGET)
	e.Post("/room/:roomid", roomPOST)
	e.Delete("/room/:roomid", roomDELETE)
	e.Get("/stream/:roomid", stream)
	if len(engine) == 0 {
		e.Run(standard.New(":8080"))
	} else {
		e.Run(fasthttp.New(":8080"))
	}
}

func stream(c han.Context) error {
	roomid := c.Param("roomid")
	listener := openListener(roomid)
	defer closeListener(roomid, listener)
	return c.SSEvent("message", listener)
}

func roomGET(c han.Context) error {
	roomid := c.Param("roomid")
	userid := fmt.Sprint(rand.Int31())
	c.Response().Header().Set(han.HeaderContentType, han.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(200)
	return html.Execute(c.Response(), han.H{
		"roomid": roomid,
		"userid": userid,
	})
}

func roomPOST(c han.Context) error {
	roomid := c.Param("roomid")
	userid := c.Form("user")
	message := c.Form("message")
	room(roomid).Submit(userid + ": " + message)

	return c.JSON(han.H{
		"status":  "success",
		"message": message,
	})
}

func roomDELETE(c han.Context) error {
	roomid := c.Param("roomid")
	deleteBroadcast(roomid)
	return nil
}
