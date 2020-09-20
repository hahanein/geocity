package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// Page is our main page component.
type Page struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (p *Page) Render() vecty.ComponentOrHTML {
	return elem.Div(
		new(Header),
		elem.Div(
			vecty.Markup(
				vecty.Class("content"),
			),
			new(Object),
			new(LoremIpsum),
		),
		new(Contact),
	)
}
