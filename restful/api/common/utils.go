package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ParsePostParam(r io.Reader, paramStruct interface{}) error {
	paramsBs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(paramsBs, paramStruct)
	if err != nil {
		return err
	}
	return nil
}
