package schema

import (
	"bytes"
	"encoding/gob"
)

func Encode(model any) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(model)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Decode(data []byte, model any) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(model)
}
