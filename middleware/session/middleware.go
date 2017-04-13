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
package session

import "github.com/admpub/han"

func Sessions(options *han.SessionOptions, store Store) han.MiddlewareFuncd {
	return func(h han.Handler) han.HandlerFunc {
		return func(c han.Context) error {
			c.SetSessionOptions(options)
			s := NewMySession(store, options.Name, c)
			c.SetSessioner(s)
			c.AddPreResponseHook(func() error {
				if options.Engine == `cookie` {
					s.Save()
				}
				return nil
			})
			err := h.Handle(c)
			s.Save()
			return err
		}
	}
}

func Middleware(options *han.SessionOptions) han.MiddlewareFuncd {
	store := StoreEngine(options)
	return Sessions(options, store)
}
