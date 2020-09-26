package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Object struct {
	vecty.Core
}

func (o *Object) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(prop.ID("object")),
	)
}
