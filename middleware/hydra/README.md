# hydra
[Hydra](https://github.com/ory-am/hydra) middleware for [Han](https://github.com/admpub/han) framework.
It uses Hydra's API to extract and validate auth token.

## Example

``` go
import (
    "github.com/admpub/han"
    "github.com/admpub/han/engine/standard"
    "github.com/ory-am/hydra/firewall"
    hydraMW "github.com/admpub/han/middleware/hydra"
)

func handler(c han.Context) error {
	ctx := c.Get("hydra").(*firewall.Context) // or hydraMW.GetContext(c)
	// Now you can access ctx.Subject etc.
	return nil
}

func main(){
	// Initialize Hydra
	hc, err := hydraMW.Connect(hydraMW.Options{
		ClientID     : "...",
		ClientSecret : "...",
		ClusterURL   : "",
	})
	if err != nil {
		panic(err)
	}

	// Use the middleware
 	e := han.New()
	e.Get("/", handler, hydraMW.ScopesRequired(hc, nil, "scope1", "scope2"))
	e.Run(standard.New(":4444"))
}
```
