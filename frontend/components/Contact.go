package components

import (
	"log"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/hahanein/geocity/entity"
	"github.com/hahanein/geocity/frontend/rest"
)

const (
	CallStateReady = iota
	CallStatePending
	CallStateDone
)

type Contact struct {
	vecty.Core
	contactState int
}

func (c *Contact) Mount() {
	c.contactState = CallStateReady

	go func() {
		c.contactState = CallStatePending
		vecty.Rerender(c)
		var reply entity.Contact
		if err := rest.Get("contact", &reply); err != nil {
			log.Panic(err)
		}
		c.contactState = CallStateDone
		vecty.Rerender(c)
	}()
}

func (f *Contact) Render() vecty.ComponentOrHTML {
	return elem.Footer(
		elem.Div(
			vecty.Markup(vecty.Class("footer__contact")),
			elem.Heading3(
				vecty.Markup(vecty.Class("footer__contact__item")),
				vecty.Text("kontakt"),
			),
			elem.Div(
				vecty.Markup(vecty.Class("footer__contact__item")),
				vecty.Text("benjamin.westphal@riseup.net"),
				elem.Break(),
				vecty.Text("+49 (0)176 20 01 38 34"),
				elem.Break(),
				elem.Break(),
				elem.Break(),
				elem.Break(),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("footer__various")),
		),
	)
}
