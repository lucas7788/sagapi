package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type ReqParam struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func ParsePostParam(r io.Reader, paramStruct interface{}) error {
	paramsBs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	rp := &ReqParam{
		Params: paramStruct,
	}
	err = json.Unmarshal(paramsBs, rp)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(rp.Params)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, paramStruct)
}
