package components

import (
	"fmt"
	"math"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Audio struct {
	vecty.Core
	title string `vecty:"props"`
	url   string `vecty:"props"`

	duration *string
}

func (a *Audio) Render() vecty.ComponentOrHTML {
	var Audio *vecty.HTML
	if a.duration == nil {
		Audio = elem.Audio(vecty.Markup(
			vecty.Class("audio"),
			prop.Src(a.url),
			&vecty.EventListener{
				Name: "canplay",
				Listener: func(e *vecty.Event) {
					total := e.Target.Get("duration").Float()
					s := math.Round(math.Mod(total, 60))
					m := math.Round(total / 60)

					text := fmt.Sprintf("%.0f:%.0f", m, s)
					a.duration = &text
					vecty.Rerender(a)
				},
			},
		))
	} else {
		Audio = elem.Audio(vecty.Markup(
			vecty.Class("audio"),
			prop.Src(a.url),
		))
	}

	var Title []vecty.MarkupOrChild
	if a.duration == nil {
		Title = []vecty.MarkupOrChild{vecty.Text(a.title)}
	} else {
		playPauseListener := &vecty.EventListener{
			Name: "click",
			Listener: func(e *vecty.Event) {
				if Audio.Node().Get("paused").Bool() {
					Audio.Node().Call("play")
				} else {
					Audio.Node().Call("pause")
				}
			},
		}
		playPauseListener.PreventDefault()

		Title = []vecty.MarkupOrChild{
			elem.Anchor(
				vecty.Markup(
					prop.Href(a.url),
					playPauseListener,
				),
				vecty.Text(a.title),
			),
			vecty.Text(fmt.Sprintf(" (%s)", *a.duration)),
		}
	}

	return elem.Figure(
		elem.FigureCaption(Title...),
		Audio,
	)
}
