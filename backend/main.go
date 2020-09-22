package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("main.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := &store{db}
	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	isFirstStart, err := store.isFirstStart()
	if err != nil {
		log.Fatal(err)
	}

	session := &session{}
	eventc := make(chan event)

	h := &handler{
		isFirstStart,
		eventc,
		session,
		store,
	}

	fmt.Println("Listening on 127.0.0.1:7001")
	log.Fatal(http.ListenAndServe("127.0.0.1:7001", h))
}

func itob(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

func btoi(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
