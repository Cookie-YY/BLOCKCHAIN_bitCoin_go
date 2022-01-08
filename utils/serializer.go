package utils

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
)

/* Serializer defines the "serialize" and "DeSerialize" func. There are two realization for now.
- 1. jsonSerializer: easy to read but goes wrong in some corner cases
- 2. gobSerializer: save the byte stream to file. It could cover most of the cases but hard to read
*/

type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	DeSerialize([]byte, interface{}) error
}

type jsonSerializer struct {
}

func NewJsonSerializer() *jsonSerializer {
	return &jsonSerializer{}
}

func (js *jsonSerializer) Serialize(content interface{}) ([]byte, error) {
	result, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (js *jsonSerializer) DeSerialize(content []byte, container interface{}) error {
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	return json.Unmarshal(content, container)
}

type gobSerializer struct {
}

func NewGobSerializer() *gobSerializer {
	return &gobSerializer{}
}

func (gs *gobSerializer) Serialize(content interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(content)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (gs *gobSerializer) DeSerialize(content []byte, container interface{}) error {
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	return decoder.Decode(container)
}
