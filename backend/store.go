package main

import (
	"github.com/boltdb/bolt"
)

type store struct {
	*bolt.DB
}

var defaultBuckets = []string{
	"inbox",
	"contact",
	"blog",
	"auth",
}

func (s *store) init() error {
	return s.Update(func(tx *bolt.Tx) error {
		for _, k := range defaultBuckets {
			_, err := tx.CreateBucketIfNotExists([]byte(k))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *store) isFirstStart() (bool, error) {
	res := false

	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("auth"))

		if bs := b.Get([]byte("username")); bs == nil {
			res = true
		}

		if bs := b.Get([]byte("hash")); bs == nil {
			res = true
		}

		return nil
	})

	return res, err
}
