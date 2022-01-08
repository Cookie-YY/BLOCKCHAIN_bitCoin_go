package utils

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strconv"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func IsExistItem(value interface{}, arr interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		arrVal := reflect.ValueOf(arr)
		for i := 0; i < arrVal.Len(); i++ {
			if reflect.DeepEqual(value, arrVal.Index(i).Interface()) {
				return true
			}
		}

	default:
		return false
	}

	return false
}

func AddIntElemIntoStringMap(stringMap map[string][]int64, mapKey string, elem int64) {
	if eleList, ok := stringMap[mapKey]; ok {
		if IsExistItem(elem, eleList) {
			return
		} else {
			eleList = append(eleList, elem)
			stringMap[mapKey] = eleList
		}
	} else {
		var eleList = make([]int64, 0)
		eleList = append(eleList, elem)
		stringMap[mapKey] = eleList
	}
}

func Str2Float(s string) (float64, error) {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func Str2Int(s string) (int64, error) {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
