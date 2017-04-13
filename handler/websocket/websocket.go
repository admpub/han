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
package websocket

import (
	"github.com/admpub/websocket"
	"github.com/admpub/han"
)

type Options struct {
	Handle   func(*websocket.Conn, han.Context) error
	Upgrader *websocket.EchoUpgrader
	Validate func(han.Context) error
	Prefix   string
}

func (o Options) Wrapper(e han.RouteRegister) {
	if o.Upgrader == nil {
		o.Upgrader = DefaultUpgrader
	}
	e.Any(o.Prefix, Websocket(o.Handle, o.Validate, o.Upgrader))
}

type Handler interface {
	Handle(*websocket.Conn, han.Context) error
	Upgrader() *websocket.EchoUpgrader
	Validate(han.Context) error
}

var (
	DefaultUpgrader = &websocket.EchoUpgrader{}
)

func HanderWrapper(v interface{}) han.Handler {
	if h, ok := v.(func(*websocket.Conn, han.Context) error); ok {
		return Websocket(h, nil)
	}
	if h, ok := v.(Handler); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader())
	}
	if h, ok := v.(Options); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader)
	}
	if h, ok := v.(*Options); ok {
		return Websocket(h.Handle, h.Validate, h.Upgrader)
	}
	if h, ok := v.(StdHandler); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader())
	}
	if h, ok := v.(StdOptions); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader)
	}
	if h, ok := v.(*StdOptions); ok {
		return StdWebsocket(h.Handle, h.Validate, h.Upgrader)
	}
	return nil
}

func Websocket(executer func(*websocket.Conn, han.Context) error, validate func(han.Context) error, opts ...*websocket.EchoUpgrader) han.HandlerFunc {
	var opt *websocket.EchoUpgrader
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = DefaultUpgrader
	}
	if executer == nil {
		//Test mode
		executer = DefaultExecuter
	}
	h := func(ctx han.Context) (err error) {
		if validate != nil {
			if err = validate(ctx); err != nil {
				return
			}
		}
		return opt.Upgrade(ctx, func(conn *websocket.Conn) error {
			defer conn.Close()
			return executer(conn, ctx)
		}, nil)
	}
	return han.HandlerFunc(h)
}
