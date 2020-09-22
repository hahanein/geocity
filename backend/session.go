package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type session struct {
	sync.RWMutex
	token  []byte
	cancel func()
}

func (s *session) Token() string {
	return base64.RawURLEncoding.EncodeToString(s.token)
}

func (s *session) Renew() error {
	s.Lock()
	defer s.Unlock()

	s.cancel()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour))
	s.cancel = cancel

	go func() {
		<-ctx.Done()
		s.token = nil
		s.cancel = func() {}
	}()

	tok := make([]byte, 32)
	if _, err := rand.Read(tok); err != nil {
		return err
	}
	s.token = tok

	return nil
}
