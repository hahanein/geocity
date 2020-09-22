package message

// TODO: Remove "optional" tag once header is automatically added.

import (
	"encoding/asn1"
)

type Any struct {
	asn1.RawContent
}

type StatusMessage = string

type Header struct {
	ID   int32
	Type string `asn1:"ia5"`
}

type AuthRequest struct {
	Username string `asn1:"ia5"`
	Password string
}

type AuthReply struct {
	StatusMessage `asn1:"optional,ia5"`
	Token         string `asn1:"optional,ia5"`
}

type PutRequest struct {
	Value  asn1.RawContent
	Bucket string `asn1:"ia5"`
	Token  string `asn1:"ia5"`
}

type PutReply struct {
	StatusMessage `asn1:"optional,ia5"`
}

type ExecRequest struct {
	Token string `asn1:"ia5"`
}

type ExecReply struct {
	StatusMessage `asn1:"optional,ia5"`
}
