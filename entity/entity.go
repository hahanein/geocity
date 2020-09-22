package entity

type Message struct {
	Email   string `asn1:"ia5"`
	Message string `asn1:"ia5"`
}

type Contact struct {
	List []string `asn1:"optional,sequence,ia5"`
}
