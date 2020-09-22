package main

import (
	"encoding/asn1"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hahanein/geocity/entity"
)

func (h *handler) restMessage(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PUT":
		var m entity.Message

		if err := request(req, &m); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		bs, err := asn1.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		if err := h.store.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("inbox"))
			t := time.Now().UnixNano()
			return b.Put(itob(t), bs)
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panic(err)
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) restContact(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		var reply entity.Contact

		if err := h.store.View(func(tx *bolt.Tx) error {
			var contact string

			b := tx.Bucket([]byte("contact"))
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				_, err := asn1.Unmarshal(v, &contact)
				if err != nil {
					return err
				}
				reply.List = append(reply.List, contact)
			}

			return nil
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		b, err := asn1.Marshal(reply)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
