package main

import (
	"bytes"
	"encoding/asn1"
	"io/ioutil"
	"log"
	"net/http"

	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/hahanein/geocity/frontend/caller"
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

func (p *App) Unmount() {
	caller.Unregister()
}

var css = `
@font-face {
	font-family: 'SctoGroteskA';
	src: url('/typeface/scto_grotesk_a_regular.woff') format('woff');
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
	width: 100%;
	padding: 0.5rem;
	background: white;
	box-sizing: border-box;
	position: sticky;
	top: 0;
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

::selection {
	background: black;
}

footer {
	border-top: 2px solid black;
	padding: 0.5rem;
	background: white;
	/* position: fixed;
	left: 0;
	bottom: 0;
	width: 100%; */
	position: sticky;
	bottom: 0;
}

footer > div, footer > h3 {
	margin-right: 1.5rem;
}

.footer_message {
	resize: none;
}

.content {
	display: flex;
}

.content > div, .content > canvas {
	flex: 1 1 auto;
}

.object {
	padding: 0.5rem;
}

.lorem-ipsum {
	border-left: 2px solid black;
	padding: 0.5rem;
	max-width: 25vw;
	min-height: 200vh;
}`

func (a *App) Mount() {
	a.init = false

	caller.Register(func(method string, params interface{}, reply interface{}) {
		b, err := asn1.Marshal(params)
		if err != nil {
			log.Panic(err)
		}

		resp, err := http.Post(
			"/api/"+method,
			"application/octet-stream",
			bytes.NewReader(b),
		)
		if err != nil {
			log.Panic(err)
		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}

		_, err = asn1.Unmarshal(b, reply)
		if err != nil {
			log.Panic(err)
		}
	})

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
