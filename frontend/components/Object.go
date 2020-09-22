package components

import (
	"log"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Object struct {
	vecty.Core
}

func (o *Object) Mount() {
	Canvas := js.Global().Get("document").Call("querySelector", "#object")
	gl := Canvas.Call("getContext", "webgl")
	if gl.IsNull() {
		log.Fatal("object: unable to initialize WebGL")
	}

	gl.Call("clearColor", 0.0, 0.0, 0.0, 0.0)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))
}

func (o *Object) Render() vecty.ComponentOrHTML {
	return elem.Canvas(
		vecty.Markup(vecty.Class("content__object"), prop.ID("object")),
	)
}
