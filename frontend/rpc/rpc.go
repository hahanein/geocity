package rpc

import (
	"bytes"
	"encoding/asn1"
	"io/ioutil"
	"net/http"
)

func Call(method string, params, reply interface{}) error {
	b, err := asn1.Marshal(params)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"/api/"+method,
		"application/octet-stream",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = asn1.Unmarshal(b, reply)
	if err != nil {
		return err
	}

	return nil
}
