package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/hahanein/geocity/frontend/caller"
	"github.com/hahanein/geocity/message"
)

const (
	CallStateReady = iota
	CallStatePending
	CallStateDone
)

type Contact struct {
	vecty.Core
	messageState int
	contactState int
}

func (c *Contact) Mount() {
	c.messageState = CallStateReady
	c.contactState = CallStateReady

	go func() {
		c.contactState = CallStatePending
		vecty.Rerender(c)
		var reply message.GetContactReply
		caller.Call("get_contact", message.GetContactRequest{}, &reply)
		c.contactState = CallStateDone
		vecty.Rerender(c)
	}()
}

func (f *Contact) Render() vecty.ComponentOrHTML {
	General := elem.Div(
		vecty.Text("benjamin.westphal@riseup.net"),
		elem.Break(),
		vecty.Text("+49 (0)176 20 01 38 34"),
		elem.Break(),
		elem.Break(),
		elem.Break(),
	)

	Email := elem.Input(vecty.Markup(
		prop.Type(prop.TypeEmail),
		prop.Name("email"),
		prop.Placeholder("email"),
	))

	Message := elem.TextArea(vecty.Markup(
		vecty.Class("footer_message"),
		prop.Name("message"),
		prop.Placeholder("nachricht (ISO/IEC 646)"),
	))

	putMessage := func(e *vecty.Event) {
		email := Email.Node().Get("value").String()
		body := Message.Node().Get("value").String()

		params := message.PutMessageRequest{
			Email:   email,
			Message: body,
		}

		go func() {
			f.messageState = CallStatePending
			vecty.Rerender(f)
			var reply message.PutMessageReply
			caller.Call("put_message", params, &reply)
			f.messageState = CallStateDone
			vecty.Rerender(f)
		}()
	}

	_ = elem.Div(
		vecty.Markup(
			vecty.Style("display", "flex"),
			vecty.Style("flex-direction", "column"),
		),
		elem.Label(
			vecty.Markup(prop.For("email")),
			vecty.Text("email"),
		),
		Email,
		elem.Break(),
		elem.Label(
			vecty.Markup(prop.For("message")),
			vecty.Text("nachricht"),
		),
		Message,
		elem.Break(),
		elem.Button(
			vecty.Markup(
				prop.Type(prop.TypeSubmit),
				event.Click(putMessage),
			),
			vecty.Text("nachricht senden"),
		),
	)

	Footer := elem.Footer(
		vecty.Markup(
			vecty.Style("display", "flex"),
		),
		elem.Heading3(vecty.Text("kontakt")),
		General,
		// Form,
	)

	switch f.messageState {
	case CallStatePending:
		return elem.Footer(
			vecty.Markup(
				vecty.Style("display", "flex"),
				vecty.Style("flex-direction", "column"),
			),
			elem.Heading3(vecty.Text("kontakt")),
			General,
			vecty.Text("Sende Nachricht..."),
		)

	case CallStateDone:
		return elem.Footer(
			vecty.Markup(
				vecty.Class("footer"),
				vecty.Style("display", "flex"),
				vecty.Style("flex-direction", "column"),
			),
			elem.Heading3(vecty.Text("kontakt")),
			General,
			vecty.Text("Nachricht gesendet."),
		)

	default:
		return Footer
	}
}
