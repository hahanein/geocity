package main

import (
	"context"
	"crypto/rand"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hahanein/geocity/message"
	"golang.org/x/crypto/bcrypt"
)

type session struct {
	Token  []byte
	Cancel func()
}

func (s *session) Renew() error {
	s.Cancel()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour))
	s.Cancel = cancel

	go func() {
		<-ctx.Done()
		s.Token = nil
		s.Cancel = func() {}
	}()

	tok := make([]byte, 32)
	if _, err := rand.Read(tok); err != nil {
		return err
	}
	s.Token = tok

	return nil
}

func main() {
	sess := session{
		Token:  nil,
		Cancel: func() {},
	}

	db, err := bolt.Open("main.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		buckets := []string{
			"inbox",
			"contact",
			"blog",
			"auth",
		}

		for _, k := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(k))
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	init := false

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("auth"))

		if bs := b.Get([]byte("username")); bs == nil {
			init = true
		}

		if bs := b.Get([]byte("hash")); bs == nil {
			init = true
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		for {
			if init {
				bs := make([]byte, 32)
				if _, err := rand.Read(bs); err != nil {
					log.Fatal(err)
				}

				id := base64.RawURLEncoding.EncodeToString(bs)

				w.Write([]byte(fmt.Sprintf("id: %s\n", id)))
				w.Write([]byte("event: init\n"))
				w.Write([]byte("data: \n\n"))
				f.Flush()
			}

			time.Sleep(5 * time.Second)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.RequestURI {
		case "/put_message":
			var m message.PutMessageRequest

			_, err := asn1.Unmarshal(bs, &m)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			if err := db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("inbox"))
				t := time.Now().UnixNano()
				return b.Put(itob(t), bs)
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			bs, err = asn1.Marshal(message.PutMessageReply{})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			break

		case "/get_contact":
			var reply message.GetContactReply

			if err := db.View(func(tx *bolt.Tx) error {
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
			break

		case "/put":
			var m message.PutRequest

			_, err := asn1.Unmarshal(bs, &m)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			if err := db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(m.Bucket))
				t := time.Now().UnixNano()
				return b.Put(itob(t), m.Value)
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			bs, err = asn1.Marshal(message.PutReply{})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			break

		case "/set_up":
			var m message.AuthRequest

			_, err := asn1.Unmarshal(bs, &m)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			if err := db.Update(func(tx *bolt.Tx) error {
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
				break
			}

			if err := sess.Renew(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			bs, err = asn1.Marshal(message.AuthReply{
				Token: base64.RawURLEncoding.EncodeToString(sess.Token),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			break

		case "/auth":
			var m message.AuthRequest

			_, err := asn1.Unmarshal(bs, &m)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			if err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("auth"))
				username := b.Get([]byte("username"))
				hash := b.Get([]byte("hash"))

				if string(username) != m.Username {
					return errors.New("unknown username")
				}

				return bcrypt.CompareHashAndPassword(hash, []byte(m.Password))
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			if err := sess.Renew(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			bs, err = asn1.Marshal(message.AuthReply{
				Token: base64.RawURLEncoding.EncodeToString(sess.Token),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			break

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	fmt.Println("Listening on 127.0.0.1:7001")
	log.Fatal(http.ListenAndServe("127.0.0.1:7001", nil))
}

func itob(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

func btoi(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
