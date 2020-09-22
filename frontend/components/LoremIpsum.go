package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type LoremIpsum struct {
	vecty.Core
}

func (li *LoremIpsum) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("content__lorem-ipsum")),
	)
}
