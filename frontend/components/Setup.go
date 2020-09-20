package components

import (
	"log"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/hahanein/geocity/frontend/caller"
	"github.com/hahanein/geocity/message"
)

// Setup is our set up component when
type Setup struct {
	vecty.Core
	state int
}

func (s *Setup) Mount() {
	s.state = CallStateReady
}

// Render implements the vecty.Component interface.
func (s *Setup) Render() vecty.ComponentOrHTML {
	Username := elem.Input(vecty.Markup(prop.Type(prop.TypeText), prop.Name("username")))
	Password := elem.Input(vecty.Markup(prop.Type(prop.TypePassword), prop.Name("password")))
	PasswordConfirmation := elem.Input(vecty.Markup(prop.Type(prop.TypePassword), prop.Name("password_confirmation")))

	setUp := func(e *vecty.Event) {
		usr := Username.Node().Get("value").String()
		pwdWant := Password.Node().Get("value").String()
		pwdHave := PasswordConfirmation.Node().Get("value").String()

		if pwdWant != pwdHave {
			log.Panic("setup: passwords do not match")
		}

		params := message.AuthRequest{
			Username: usr,
			Password: pwdWant,
		}

		go func() {
			s.state = CallStatePending
			vecty.Rerender(s)
			var reply message.AuthReply
			caller.Call("set_up", params, &reply)
			s.state = CallStateDone
			vecty.Rerender(s)
		}()
	}

	Form := elem.Div(
		vecty.Markup(
			vecty.Style("display", "flex"),
			vecty.Style("flex-direction", "column"),
		),
		elem.Heading3(vecty.Text("Setup")),
		elem.Label(
			vecty.Markup(prop.For("username")),
			vecty.Text("Username"),
		),
		Username,
		elem.Break(),
		elem.Label(
			vecty.Markup(prop.For("password")),
			vecty.Text("Password"),
		),
		Password,
		elem.Break(),
		elem.Label(
			vecty.Markup(prop.For("password_confirmation")),
			vecty.Text("Password confirmation"),
		),
		PasswordConfirmation,
		elem.Break(),
		elem.Button(
			vecty.Markup(
				prop.Type(prop.TypeSubmit),
				event.Click(setUp),
			),
			vecty.Text("Submit"),
		),
	)
	return Form
}
