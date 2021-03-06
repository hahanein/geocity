package main

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/hahanein/geocity/frontend/components"
)

func main() {
	vecty.SetTitle("Benjamin Westphal")
	vecty.RenderBody(new(App))
}

type App struct {
	vecty.Core
	init bool
}

var css = `
@font-face {
	font-family: 'SctoGroteskA';
	src: url('/typeface/scto_grotesk_a_regular.woff') format('woff');
	font-weight: normal;
	font-style: normal;
}

* {
	font-weight: normal;
	font-size: 1.1rem;
	margin: 0;
	padding: 0;
}

body {
	font-family: 'SctoGroteskA', sans-serif;
	-webkit-font-smoothing: antialiased;
}

button {
	background: none;
	border: none;
}

label {
	display: none;
}

input, button {
	font-size: 0.97rem;
}

br {
	margin-top: 0;
}

a {
	color: black;
}

ol {
	padding-left: 1.2em;
}

::selection {
	background: black;
}

#root {
	box-sizing: border-box;
	padding: 0.55rem;
	width: 100vw;
	height: 100vh;
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(min(5.5rem, 100%), 1fr));
	grid-template-rows: repeat(auto-fill, minmax(min(5.5rem, 100%), 1fr));
	grid-gap: 0.55rem;
}

#title {
	grid-row-start: 1;
	grid-column-start: 1;
	grid-column-end: span2;
}

@media only screen and (max-width: 767px) {
	#route__current {
		grid-row-start: 1;
		grid-column-start: 3;
		grid-column-end: span 4;
	}

	#stream {
		grid-row-start: -3;
		grid-column-start: 1;
		grid-column-end: span 4;
		padding-bottom: 1.1rem;
	}
}

@media only screen and (min-width: 768px) {
	#route__current {
		grid-row-start: 1;
		grid-column-start: -5;
		grid-column-end: span 4;
	}

	#stream {
		grid-row-start: -3;
		grid-column-start: -5;
		grid-column-end: span 4;
	}
}

#contact__heading {
	grid-column-start: 1;
	grid-row-start: -2;
}

#contact__info {
	grid-column-start: 2;
	grid-column-end: span 3;
	grid-row-start: -2;
}

#object {
	grid-row-start: 2;
}

#stream {
	display: flex;
	flex-direction: column;
	justify-content: flex-end;
}`

func (a *App) Mount() {
	a.init = false

	Style := js.Global().Get("document").Call("createElement", "style")
	Style.Set("innerHTML", css)

	Viewport := js.Global().Get("document").Call("createElement", "meta")
	Viewport.Set("name", "viewport")
	Viewport.Set("content", "width=device-width, initial-scale=1")

	js.Global().Get("document").Call("getElementsByTagName", "head").Index(0).Call("appendChild", Style)
	js.Global().Get("document").Call("getElementsByTagName", "head").Index(0).Call("appendChild", Viewport)

	evtSource := js.Global().Get("EventSource").New("/api/events")
	evtSource.Call("addEventListener", "init", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		a.init = true
		vecty.Rerender(a)

		args[0].Get("target").Call("removeEventListener", "init", this)
		return nil
	}))
}

func (a *App) Render() vecty.ComponentOrHTML {
	if a.init {
		return elem.Body(new(components.Setup))
	}

	return elem.Body(new(components.Page))
}
