package main

import (
	"encoding/asn1"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hahanein/geocity/message"
	"golang.org/x/crypto/bcrypt"
)

func request(req *http.Request, v interface{}) error {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	_, err = asn1.Unmarshal(bs, v)
	if err != nil {
		return err
	}

	return nil
}

func reply(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/octet-stream")

	bs, err := asn1.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (h *handler) rpcPut(w http.ResponseWriter, req *http.Request) {
	var m message.PutRequest
	if err := request(req, &m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	if err := h.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(m.Bucket))
		t := time.Now().UnixNano()
		return b.Put(itob(t), m.Value)
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	reply(w, message.PutReply{})
}

func (h *handler) rpcSetUp(w http.ResponseWriter, req *http.Request) {
	var m message.AuthRequest
	if err := request(req, &m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	if err := h.store.Update(func(tx *bolt.Tx) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte("auth"))
		b.Put([]byte("username"), []byte(m.Username))
		b.Put([]byte("hash"), hash)

		return nil
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	if err := h.session.Renew(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	reply(w, message.AuthReply{Token: h.session.Token()})
}

func (h *handler) rpcAuthorize(w http.ResponseWriter, req *http.Request) {
	var m message.AuthRequest
	if err := request(req, &m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	if err := h.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("auth"))
		username := b.Get([]byte("username"))
		hash := b.Get([]byte("hash"))

		if string(username) != m.Username {
			return errors.New("unknown username")
		}

		return bcrypt.CompareHashAndPassword(hash, []byte(m.Password))
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	if err := h.session.Renew(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	reply(w, message.AuthReply{Token: h.session.Token()})
}
