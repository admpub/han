# Han
[![Build Status](https://travis-ci.org/admpub/han.svg?branch=master)](https://travis-ci.org/admpub/han) [![Go Report Card](https://goreportcard.com/badge/github.com/admpub/han)](https://goreportcard.com/report/github.com/admpub/han)
#### Han 是一个快速而简洁的Web通用开发框架
使用本框架需要 Go 版本 >= **go 1.7**

## 框架特色

- 经过优化的HTTP路由
- 容易构建RESTful API服务
- 支持原生HTTP和FastHTTP
- 支持路由分组
- 支持用中间件扩展框架
- 内置常用的HTTP响应类型支持
- 支持任意模板引擎
- 内置Session、数据验证和多语言接口
- 提供便捷的Session和Cookie操作方法
- 丰富的中间件

## 快速入门

### 安装

```sh
$ go get github.com/admpub/han
```

### Hello, World!

创建文件 `server.go`

```go
package main

import (
	"net/http"
	"github.com/admpub/han"
	"github.com/admpub/han/engine/standard"
)

func main() {
	e := han.New()
	e.Get("/", func(c han.Context) error {
		return c.String("Hello, World!", http.StatusOK)
	})
	e.Run(standard.New(":1323"))
}
```

运行 server

```sh
$ go run server.go
```

用浏览器打开 [http://localhost:1323](http://localhost:1323) ，在页面上你将看到
`Hello, World!`。

### 路由

```go
e.Post("/users", saveUser)
e.Get("/users/:id", getUser)
e.Put("/users/:id", updateUser)
e.Delete("/users/:id", deleteUser)
```

### 网址路径参数

```go
func getUser(c han.Context) error {
	// 从`users/:id`中获取id参数值
	id := c.Param("id")
}
```

### 网址查询参数(Query Parameters)

`/show?team=x-men&member=wolverine`

```go
func show(c han.Context) error {
	// 从查询字符串中获取team和member值
	team := c.Query("team")
	member := c.Query("member")
}
```

### 表单 `application/x-www-form-urlencoded`

`POST` `/save`

name | value
:--- | :---
name | Joe Smith
email | joe@labstack.com


```go
func save(c han.Context) error {
	// 从表单中获取 name 和 email
	name := c.Form("name")
	email := c.Form("email")
}
```

### 表单 `multipart/form-data`

`POST` `/save`

name | value
:--- | :---
name | Joe Smith
email | joe@labstack.com
avatar | avatar

```go
func save(c han.Context) error {
	// 从表单中获取 name 和 email
	name := c.Form("name")
	email := c.Form("email")

	//------------
	// 获取附件 avatar
	//------------
	_, err := c.SaveUploadedFile("avatar","./")
	return err
}
```

### 处理请求

- Bind `JSON` or `XML` payload into Go struct based on `Content-Type` request header.
- Render response as `JSON` or `XML` with status code.

```go
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

e.Post("/users", func(c han.Context) error {
	u := new(User)
	if err := c.MustBind(u); err != nil {
		return err
	}
	return c.JSON(u, http.StatusCreated)
	// or
	// return c.XML(u, http.StatusCreated)
})
```

### 静态内容

Server any file from static directory for path `/static/*`.

```go
e.Use(mw.Static(&mw.StaticOptions{
	Root:"static", //存放静态文件的物理路径
	Path:"/static/", //网址访问静态文件的路径
	Browse:true, //是否在首页显示文件列表
}))
```

### 使用中间件

```go
// Root level middleware
e.Use(middleware.Log())
e.Use(middleware.Recover())

// Group level middleware
g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string) bool {
	if username == "joe" && password == "secret" {
		return true
	}
	return false
}))

// Route level middleware
track := func(next han.HandlerFunc) han.HandlerFunc {
	return func(c han.Context) error {
		println("request to /users")
		return next.Handle(c)
	}
}
e.Get("/users", func(c han.Context) error {
	return c.String("/users", http.StatusOK)
}, track)
```

### Cookie
```go
e.Get("/setcookie", func(c han.Context) error {
	c.SetCookie("uid","1")
	return c.String("/setcookie: uid="+c.GetCookie("uid"), http.StatusOK)
})
```

### Session
```go
...
import (
	...
	"github.com/admpub/han/middleware/session"
	//boltStore "github.com/admpub/han/middleware/session/engine/bolt"
	cookieStore "github.com/admpub/han/middleware/session/engine/cookie"
)
...
sessionOptions := &han.SessionOptions{
	Engine: `cookie`,
	Name:   `SESSIONID`,
	CookieOptions: &han.CookieOptions{
		Path:     `/`,
		Domain:   ``,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	},
}

cookieStore.RegWithOptions(&cookieStore.CookieOptions{
	KeyPairs: [][]byte{
		[]byte(`123456789012345678901234567890ab`),
	},
	SessionOptions: sessionOptions,
})

e.Use(session.Middleware(sessionOptions))

e.Get("/session", func(c han.Context) error {
	c.Session().Set("uid",1).Save()
	return c.String(fmt.Sprintf("/session: uid=%v",c.Session().Get("uid")))
})
```

### Websocket
```go
...
import (
	...
	"github.com/admpub/websocket"
	"github.com/admpub/han"
	ws "github.com/admpub/han/handler/websocket"
)
...

e.AddHandlerWrapper(ws.HanderWrapper)

e.Get("/websocket", func(c *websocket.Conn, ctx han.Context) error {
	//push(writer)
	go func() {
		var counter int
		for {
			if counter >= 10 { //测试只推10条
				return
			}
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			ctx.Logger().Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				ctx.Logger().Error(`Push error: `, err.Error())
				return
			}
			counter++
		}
	}()

	//han
	ws.DefaultExecuter(c, ctx)
	return nil
})
```
[More...](https://github.com/admpub/han/blob/master/handler/websocket/example/main.go)

### Sockjs
```go
...
import (
	...
	"github.com/admpub/han"
	"github.com/admpub/sockjs-go/sockjs"
	ws "github.com/admpub/han/handler/sockjs"
)
...

options := ws.Options{
	Handle: func(c sockjs.Session) error {
		//push(writer)
		go func() {
			var counter int
			for {
				if counter >= 10 { //测试只推10条
					return
				}
				time.Sleep(5 * time.Second)
				message := time.Now().String()
				log.Info(`Push message: `, message)
				if err := c.Send(message); err != nil {
					log.Error(`Push error: `, err.Error())
					return
				}
				counter++
			}
		}()

		//han
		ws.DefaultExecuter(c)
		return nil
	},
	Options: &sockjs.DefaultOptions,
	Prefix:  "/websocket",
}
options.Wrapper(e)
```
[More...](https://github.com/admpub/han/blob/master/handler/sockjs/example/main.go)

### Other Example

```go
package main

import (
	"net/http"

	"github.com/admpub/han"
	// "github.com/admpub/han/engine/fasthttp"
	"github.com/admpub/han/engine/standard"
	mw "github.com/admpub/han/middleware"
)

func main() {
	e := han.New()
	e.Use(mw.Log())

	e.Get("/", func(c han.Context) error {
		return c.String("Hello, World!")
	})
	e.Get("/han/:name", func(c han.Context) error {
		return c.String("Han " + c.Param("name"))
	})
	
	e.Get("/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`standard net/http handleFunc`))
		w.WriteHeader(200)
	})

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
```

[See other examples...](https://github.com/admpub/han-example/blob/master/_v2/main.go)

## 内置中间件列表
中间件  | 导入路径 | 说明
:-----------|:------------|:-----------
[BasicAuth](https://github.com/admpub/han/blob/master/middleware/auth.go)  | github.com/admpub/han/middleware |HTTP basic authentication
[BodyLimit](https://github.com/admpub/han/blob/master/middleware/bodylimit.go)  | github.com/admpub/han/middleware |Limit request body
[Gzip](https://github.com/admpub/han/blob/master/middleware/compress.go)  | github.com/admpub/han/middleware |Send gzip HTTP response
[Secure](https://github.com/admpub/han/blob/master/middleware/secure.go)  | github.com/admpub/han/middleware |Protection against attacks
[CORS](https://github.com/admpub/han/blob/master/middleware/cors.go)  | github.com/admpub/han/middleware |Cross-Origin Resource Sharing
[CSRF](https://github.com/admpub/han/blob/master/middleware/csrf.go)  | github.com/admpub/han/middleware |Cross-Site Request Forgery
[Log](https://github.com/admpub/han/blob/master/middleware/log.go)  | github.com/admpub/han/middleware |Log HTTP requests
[MethodOverride](https://github.com/admpub/han/blob/master/middleware/methodOverride.go)  | github.com/admpub/han/middleware |Override request method
[Recover](https://github.com/admpub/han/blob/master/middleware/recover.go)  | github.com/admpub/han/middleware |Recover from panics
[HTTPSRedirect](https://github.com/admpub/han/blob/master/middleware/redirect.go)  | github.com/admpub/han/middleware |Redirect HTTP requests to HTTPS
[HTTPSWWWRedirect](https://github.com/admpub/han/blob/master/middleware/redirect.go)  | github.com/admpub/han/middleware |Redirect HTTP requests to WWW HTTPS
[WWWRedirect](https://github.com/admpub/han/blob/master/middleware/redirect.go)  | github.com/admpub/han/middleware |Redirect non WWW requests to WWW
[NonWWWRedirect](https://github.com/admpub/han/blob/master/middleware/redirect.go)  | github.com/admpub/han/middleware |Redirect WWW requests to non WWW
[AddTrailingSlash](https://github.com/admpub/han/blob/master/middleware/slash.go)  | github.com/admpub/han/middleware |Add trailing slash to the request URI
[RemoveTrailingSlash](https://github.com/admpub/han/blob/master/middleware/slash.go)  | github.com/admpub/han/middleware |Remove trailing slash from the request URI
[Static](https://github.com/admpub/han/blob/master/middleware/static.go)  | github.com/admpub/han/middleware |Serve static files
[MaxAllowed](https://github.com/admpub/han/blob/master/middleware/limit.go) | github.com/admpub/han/middleware | MaxAllowed limits simultaneous requests; can help with high traffic load
[RateLimit](https://github.com/admpub/han/tree/master/middleware/ratelimit) | github.com/admpub/han/middleware/ratelimit | Rate limiting HTTP requests
[Language](https://github.com/admpub/han/tree/master/middleware/language) | github.com/admpub/han/middleware/language | Multi-language support
[Session](https://github.com/admpub/han/blob/master/middleware/session/middleware.go)  | github.com/admpub/han/middleware/session | Sessions Manager
[JWT](https://github.com/admpub/han/blob/master/middleware/jwt/jwt.go)  | github.com/admpub/han/middleware/jwt | JWT authentication
[Hydra](https://github.com/admpub/han/blob/master/middleware/hydra/hydra.go)  | github.com/admpub/han/middleware/hydra | It uses [Hydra](https://github.com/ory-am/hydra)'s API to extract and validate auth token.
[Markdown](https://github.com/admpub/han/blob/master/middleware/markdown/markdown.go)  | github.com/admpub/han/middleware/markdown | Markdown rendering
[Render](https://github.com/admpub/han/blob/master/middleware/render/middleware.go)  | github.com/admpub/han/middleware/render | HTML template rendering
[ReverseProxy](https://github.com/webx-top/reverseproxy/blob/master/middleware.go)  | github.com/webx-top/reverseproxy | Reverse proxy


## 内置处理器
功能     | 导入路径 | 说明
:-----------|:------------|:-----------
Websocket   |github.com/admpub/han/handler/websocket | [Example](https://github.com/admpub/han/blob/master/handler/websocket/example/main.go)
Sockjs      |github.com/admpub/han/handler/sockjs | [Example](https://github.com/admpub/han/blob/master/handler/sockjs/example/main.go)
Oauth2      |github.com/admpub/han/handler/oauth2 | [Example](https://github.com/admpub/han/blob/master/handler/oauth2/example/main.go)
Pprof      |github.com/admpub/han/handler/pprof | -
MVC      |github.com/admpub/han/handler/mvc | [Example](https://github.com/admpub/han/blob/master/handler/mvc/test/main.go)

## License

[Apache 2](https://github.com/admpub/han/blob/master/LICENSE)
