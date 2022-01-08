package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
)

type jsonSerializer struct {
}

func NewJsonSerializer() *jsonSerializer {
	return &jsonSerializer{}
}

func (js *jsonSerializer) Serialize(content interface{}) []byte {
	result, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (js *jsonSerializer) DeSerialize(content []byte, container interface{}) {
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))
	err := json.Unmarshal(content, container)
	if err != nil {
		log.Fatal(err)
	}
}

type gobSerializer struct {
}

func NewGobSerializer() *gobSerializer {
	return &gobSerializer{}
}

func (gs *gobSerializer) Serialize(content interface{}) []byte {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(content)
	if err != nil {
		log.Fatal(err)
	}
	return buffer.Bytes()
}

func (gs *gobSerializer) DeSerialize(content []byte, container interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err := decoder.Decode(container)
	if err != nil {
		log.Fatal(err)
	}
}

type Person struct {
	Name string
}

func main() {
	p := Person{"a"}
	// json
	resJson := NewJsonSerializer().Serialize(p)
	containerJson := Person{}
	NewGobSerializer().DeSerialize(resJson, &containerJson)
	fmt.Printf("gob serialize: %v", resJson)
	fmt.Printf("gob deSerialize: %v", containerJson)

	// gob
	resGob := NewGobSerializer().Serialize(p)
	containerGob := Person{}
	NewGobSerializer().DeSerialize(resGob, &containerGob)
	fmt.Printf("gob serialize: %v", resGob)
	fmt.Printf("gob deSerialize: %v", containerGob)
}
