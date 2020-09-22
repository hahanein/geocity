package main

import "encoding/asn1"

type event struct {
	ID    string
	Name  string
	Data  asn1.RawContent
	Retry int
}
