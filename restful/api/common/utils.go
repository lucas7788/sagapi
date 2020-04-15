package common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type ReqParam struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func ParsePostParam(c *gin.Context, paramStruct interface{}) error {
	ontid, ok := c.Get("Ontid")
	if !ok || ontid == nil {
		return fmt.Errorf("ontid is nil")
	}
	fmt.Println("ontid:", ontid)
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
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
