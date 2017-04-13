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
package mvc

import (
	"github.com/admpub/han"
	"github.com/admpub/han/logger"
)

func NewController(c han.Context) *Controller {
	a := &Controller{}
	a.Init(c)
	return a
}

type Controller struct {
	*Context
	logger.Logger
}

func (a *Controller) Init(c han.Context) error {
	a.Context = c.(*Context)
	a.Logger = c.Han().Logger()
	return nil
}
