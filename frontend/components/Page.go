package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// Page is our main page component.
type Page struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (p *Page) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			prop.ID("root"),
		),
		elem.Heading1(
			vecty.Markup(prop.ID("title")),
			vecty.Text("Benjamin Westphal"),
		),
		elem.Div(
			vecty.Markup(prop.ID("route__current")),
			vecty.Text("Bild und Ton"),
		),
		new(Object),
		new(Stream),
		elem.Heading3(
			vecty.Markup(prop.ID("contact__heading")),
			vecty.Text("Kontakt"),
		),
		elem.Div(
			vecty.Markup(prop.ID("contact__info")),
			vecty.Text("benjamin.westphal@riseup.net"),
			elem.Break(),
			vecty.Text("+49 (0)176 20 01 38 34"),
			elem.Break(),
		),
	)
}
