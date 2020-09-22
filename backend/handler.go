package main

import (
	"net/http"
)

type handler struct {
	isFirstStart bool
	eventc       chan event
	session      *session
	store        *store
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	// SSE
	case "/events":
		h.sseEvents(w, req)

	// REST
	case "/message":
		h.restMessage(w, req)
	case "/contact":
		h.restContact(w, req)

	// RPC
	case "/put":
		h.rpcPut(w, req)
	case "/set_up":
		h.rpcSetUp(w, req)
	case "/authorize":
		h.rpcAuthorize(w, req)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
