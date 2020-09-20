package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Header struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (h *Header) Render() vecty.ComponentOrHTML {
	return elem.Header(
		vecty.Markup(
			vecty.Class("header"),
		),
		elem.Heading1(
			vecty.Text("benjamin westphal "),
			// elem.Image(vecty.Markup(
			// 	vecty.Style("max-height", "1.1rem"),
			// 	prop.Src("/media/period.svg")),
			// ),
		),
	)
}
