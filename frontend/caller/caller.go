package caller

import (
	"log"
)

func none(_ string, _, _ interface{}) {
	log.Fatal("caller: no callback registered")
}

var callback = none

func Call(method string, params, reply interface{}) {
	callback(method, params, reply)
}

func Register(fn func(string, interface{}, interface{})) {
	callback = fn
}

func Unregister() {
	callback = none
}
