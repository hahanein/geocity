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

type Footer struct {
	vecty.Core
	contactState int
}

func (f *Footer) Mount() {
	f.contactState = CallStateReady

	go func() {
		f.contactState = CallStatePending
		vecty.Rerender(f)
		var reply entity.Contact
		if err := rest.Get("contact", &reply); err != nil {
			log.Panic(err)
		}
		f.contactState = CallStateDone
		vecty.Rerender(f)
	}()
}

func (f *Footer) Render() vecty.ComponentOrHTML {
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
