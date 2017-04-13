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
package pprof

import (
	"github.com/admpub/han"
)

// Wrap adds several routes from package `net/http/pprof` to *gin.Engine object
func Wrap(router *han.Han) {
	router.Get("/debug/pprof/", IndexHandler())
	router.Get("/debug/pprof/heap", HeapHandler())
	router.Get("/debug/pprof/goroutine", GoroutineHandler())
	router.Get("/debug/pprof/block", BlockHandler())
	router.Get("/debug/pprof/threadcreate", ThreadCreateHandler())
	router.Get("/debug/pprof/cmdline", CmdlineHandler())
	router.Get("/debug/pprof/profile", ProfileHandler())
	router.Get("/debug/pprof/symbol", SymbolHandler())
	router.Get("/debug/pprof/trace", TraceHandler())
}

// Wrapper make sure we are backward compatible
var Wrapper = Wrap

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Index(ctx)
		return nil
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Handler("heap").ServeHTTP(ctx)
		return nil
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Handler("goroutine").ServeHTTP(ctx)
		return nil
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Handler("block").ServeHTTP(ctx)
		return nil
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Handler("threadcreate").ServeHTTP(ctx)
		return nil
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Cmdline(ctx)
		return nil
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Profile(ctx)
		return nil
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Symbol(ctx)
		return nil
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler() han.HandlerFunc {
	return func(ctx han.Context) error {
		Trace(ctx)
		return nil
	}
}
