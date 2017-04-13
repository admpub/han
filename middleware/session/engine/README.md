# Session

Middleware support for han, utilizing by
[admpub/sessions](https://github.com/admpub/sessions).

## Installation

```shell
go get github.com/admpub/han
```

## Usage

```go
package main

import (
    "github.com/admpub/han"
    "github.com/admpub/han/engine/standard"
    "github.com/admpub/han/middleware/session"
    cookieStore "github.com/admpub/han/middleware/session/engine/cookie"
)

func index(c han.Context) error {
    session := c.Session()

    var count int
    v := session.Get("count")

    if v == nil {
        count = 0
    } else {
        count = v.(int)
        count += 1
    }

    session.Set("count", count)

    data := struct {
        Visit int
    }{
        Visit: count,
    }

    return c.JSON(http.StatusOK, data)
}

func main() {
    sessionOptions:=&han.SessionOptions{
        Name:   `GOSESSIONID`,
        Engine: `cookie`,
        CookieOptions: &han.CookieOptions{
            Path:     `/`,
            HttpOnly: true,
        },
    }
    cookieStore.RegWithOptions(&cookieStore.CookieOptions{
        KeyPairs: [][]byte{
            []byte("secret-key"),
        },
        SessionOptions: sessionOptions,
    })

    e := han.New()

    // Attach middleware
    e.Use(session.Middleware(sessionOptions))

    // Routes
    e.Get("/", index)

    e.Run(standard.New(":8080"))
}
```
