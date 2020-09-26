package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

type Stream struct {
	vecty.Core
}

func (s *Stream) Render() vecty.ComponentOrHTML {
	var Items []vecty.MarkupOrChild

	for _, url := range []string{
		"https://github.com/hahanein/ep1/blob/master/01_UNBENANNT_0280.mp3?raw=true",
		"https://github.com/hahanein/ep1/blob/master/02_UNBENANNT_0323.mp3?raw=true",
		"https://github.com/hahanein/ep1/blob/master/03_UNBENANNT_0346.mp3?raw=true",
		"https://github.com/hahanein/ep1/blob/master/04_UNBENANNT_0253.mp3?raw=true",
		"https://github.com/hahanein/ep1/blob/master/05_UNBENANNT_0315.mp3?raw=true",
	} {
		Items = append(Items, elem.ListItem(&Audio{title: "Unbenannt", url: url}))
	}

	List := elem.OrderedList(Items...)

	return elem.Div(
		vecty.Markup(prop.ID("stream")),
		elem.Div(
			vecty.Markup(vecty.Class("stream__item")),
			elem.Div(
				elem.Heading2(vecty.Text("Unbenannt")),
				elem.Heading3(
					vecty.Text("mit "),
					elem.Anchor(
						vecty.Markup(prop.Href("https://ruhestoerung.noblogs.org/")),
						vecty.Text("Notorische Ruhestörung"),
					),
				),
				elem.Break(),
				elem.Paragraph(vecty.Text("Aufgenommen irgendwo in Berlin vom 19.5.19 bis zur Unterbrechung.")),
			),
			elem.Break(),
			List,
			elem.Div(
				elem.Break(),
				vecty.Text("Notorische Ruhestörung: Stimme, Schlagzeug"),
				elem.Break(),
				vecty.Text("Benjamin: Stimme, Gitarre, Weckglas"),
			),
		),
	)
}
