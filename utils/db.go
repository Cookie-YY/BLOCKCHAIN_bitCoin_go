package utils

import (
	"io/ioutil"
	"os"
)

// DB : defines read and write func to persist block-chain
type DB interface {
	ReadFromDB(interface{}) error
	WriteToDB(interface{}) error
}

type JsonDB struct {
	path       string
	serializer Serializer
	DB
}

func GetJsonDB(path string, serializer Serializer) *JsonDB {
	return &JsonDB{path: path, serializer: serializer}
}

func (jdb *JsonDB) IsExist() bool {
	_, err := os.Stat(jdb.path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (jdb *JsonDB) ReadFromDB(contentStruct interface{}) error {
	// open file
	f, err := os.Open(jdb.path)
	if err != nil {
		return err
	}
	// read content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	// parse content
	err = jdb.serializer.DeSerialize(content, contentStruct)
	if err != nil {
		return err
	}
	// close file
	return f.Close()
}

func (jdb *JsonDB) WriteToDB(content interface{}) error {
	// format content
	result, err := jdb.serializer.Serialize(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jdb.path, result, 0644) // ignore_security_alert: there's no way to change path externally
}
