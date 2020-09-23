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
	/* background: #FFF7C7; */
	-webkit-font-smoothing: antialiased;
}

h1 {
	/* font-size: 1.8rem; */
}

header {
	border-bottom: 2px solid black;
	position: sticky;
	display: flex;
}

.header__title {
	padding: 0.5rem;
	flex: 1 1 auto;
}

.header__navigation {
	border-left: 2px solid black;
	width: calc(25vw - 1rem);
	padding: 0.5rem;
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

footer {
	border-top: 2px solid black;
	background: white;
	position: sticky;
	bottom: 0;
	display: flex;
	justify-content: space-between;
}

.footer__contact {
	display: flex;
	padding: 0.5rem;
}

.footer__contact__item {
	margin-right: 1.5rem;
	flex: 1 1 auto;
}

.footer__various {
	border-left: 2px solid black;
	width: 25vw;
}

.footer_message {
	resize: none;
}

.content {
	display: flex;
}

.content__object {
	padding: 0.5rem;
	flex: 5 5 auto;
}

.content__stream {
	border-left: 2px solid black;
	width: 25vw;
	height: calc(100vh - 10.4rem);
	display: flex;
	flex-direction: column;
	justify-content: flex-end;
}

.content__stream__item {
	padding: 0.5rem;
	border-top: 2px solid black;
}`

func (a *App) Mount() {
	a.init = false

	Style := js.Global().Get("document").Call("createElement", "style")
	Style.Set("innerHTML", css)
	js.Global().Get("document").Call("getElementsByTagName", "head").Index(0).Call("appendChild", Style)

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
