package main

import (
	"html/template"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/admpub/han"
	"github.com/admpub/han/engine/fasthttp"
	"github.com/admpub/han/engine/standard"
	ws "github.com/admpub/han/handler/sockjs"
	mw "github.com/admpub/han/middleware"
)

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
</head>
<body>
Chat : <input id='chat' type='text' value='' size='20'/>
<button onclick="sent()">Sent</button>
<div id='content'></div>

<script src="{{.sockjs}}"></script>
<script>
function $D(tag){
	return document.createElement(tag||'div');
}
function $A(ele, father){
	father = father || document.body;
	father.appendChild(ele);
}
function $(tag){
	return document.getElementById(tag);
}
function I(ele, str){
	ele.innerHTML = str||'';
}

function showMsg(msg){
	I($('content'), $('content').innerHTML + '<br >' + msg);

}
function sent(){
	ws.send($('chat').value);
}

var ws = new SockJS('{{.han}}');
ws.onopen    = function(){
	showMsg('onopen');
};
ws.onclose   = function(){
	showMsg('onclose');
};
ws.onmessage = function(msg){
	showMsg(msg.data);
};

</script>
</body>
</html>
`))

func main() {
	e := han.New()
	e.Use(mw.Log())

	e.Get("/", func(c han.Context) error {
		homeTemplate.Execute(c.Response(), map[string]string{
			"han":   "http://" + c.Request().Host() + "/websocket",
			"sockjs": sockjs.DefaultOptions.SockJSURL,
			"notice": "http://" + c.Request().Host() + "/notice",
		})
		return nil
	})

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
			// */

			//han
			ws.DefaultExecuter(c)
			return nil
		},
		Options: &sockjs.DefaultOptions,
		Prefix:  "/websocket",
	}
	options.Wrapper(e)

	options.Handle = func(c sockjs.Session) error {
		for {
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			log.Info(`Push message: `, message)
			if err := c.Send(message); err != nil {
				return err
			}
		}
		return nil
	}
	options.Prefix = "/notice"

	options.Wrapper(e)

	switch `` {
	case `fast`:
		panic(`Unimplemented`)
		// FastHTTP
		e.Run(fasthttp.New(":4444"))

	default:
		// Standard
		e.Run(standard.New(":4444"))
	}
}
