package rest

import (
	"bytes"
	"encoding/asn1"
	"io/ioutil"
	"net/http"
)

func Put(resource string, v interface{}) error {
	b, err := asn1.Marshal(v)
	if err != nil {
		return err
	}

	_, err = http.NewRequest(
		http.MethodPut,
		"/api/"+resource,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	return nil
}

func Get(resource string, v interface{}) error {
	resp, err := http.Get("/api/" + resource)
	if err != nil {
		return err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = asn1.Unmarshal(bs, v)
	if err != nil {
		return err
	}

	return nil
}
