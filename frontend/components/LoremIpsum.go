package components

import (
	"math/rand"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type LoremIpsum struct {
	vecty.Core
}

var quotes = []string{
	`So ließen sich vielleicht aus allen Bauern und Handwerkern Künstler bilden, d. h. Menschen, die ihr Gewerbe um ihres Gewerbes willen liebten, durch eigengelenkte Kraft und eigne Erfindsamkeit verbesserten und dadurch ihre intellektuellen Kräfte kultivierten, ihren Charakter vereitelten, ihre Genüsse erhöhten. So würde die Menschheit durch eben die Dinge geadelt, die jetzt, wie schön sie auch an sich sind, so oft dazu dienen, sie zu entehren. Je mehr der Mensch in Ideen und Empfindungen zu leben gewohnt ist, je stärker und feiner seine intellektuelle und moralische Kraft ist, desto mehr sucht er allein solche äußre Lagen zu wählen, welche zugleich dem innren Menschen mehr Stoff geben, oder denjenigen, in welche ihn das Schicksal wirft, wenigstens solche Seiten abzugewinnen. Der Gewinn, welchen der Mensch an Größe und Schönheit einerntet, wenn er unaufhörlich dahin strebt, daß sein inneres Dasein immer den ersten Platz behaupte, daß es immer der erste Quell und das letzte Ziel alles Wirkens und alles Körperliche und Äußere nur Hülle und Werkzeug desselben sei, ist unabsehlich.`,
	`Allein, freilich ist Freiheit die notwendige Bedingung, ohne welche selbst das seelenvollste Geschäft keine heilsamen Wirkungen dieser Art hervorzubringen vermag. Was nicht von dem Menschen selbst gewählt, worin er auch nur eingeschränkt und geleitet wird, das geht nicht in sein Wesen über, das bleibt ihm ewig fremd, das verrichtet er nicht eigentlich mit menschlicher Kraft, sondern mit mechanischer Fertigkeit.`,
}

func (li *LoremIpsum) Render() vecty.ComponentOrHTML {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(quotes))

	return elem.Div(
		vecty.Markup(vecty.Class("lorem-ipsum")),
		elem.Paragraph(vecty.Text(quotes[i])),
		elem.Break(),
		elem.Paragraph(
			vecty.Text(`— Wilhelm von Humboldt`),
			vecty.Markup(vecty.Style("text-align", "right")),
		),
	)
}
