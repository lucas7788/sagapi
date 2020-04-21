package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseResult struct {
	Result  interface{} `json:"result"`
	Error   int64       `json:"error"`
	Desc    string      `json:"description"`
	Version string      `json:"version"`
}

func WriteResponse(c *gin.Context, response *ResponseResult) error {
	bs, err := json.Marshal(response)
	if err != nil {
		return err
	}

	c.String(http.StatusOK, string(bs))
	return nil
}

func ResponseSuccess(result interface{}) *ResponseResult {
	return &ResponseResult{
		Result:  result,
		Error:   SUCCESS,
		Desc:    ErrMap[SUCCESS],
		Version: "1.0",
	}
}

func ResponseFailed(errCode int64, err error) *ResponseResult {
	return &ResponseResult{
		Result:  nil,
		Error:   errCode,
		Desc:    ErrMap[errCode] + err.Error(),
		Version: "1.0",
	}
}

var ErrMap = map[int64]string{
	SUCCESS:           "SUCCESS",
	PARA_ERROR:        "PARAMETER ERROR",
	INTER_ERROR:       "INTER_ERROR",
	SQL_ERROR:         "SQL_ERROR",
	JWT_VERIFY_FAILED: "JWT_VERIFY_FAILED",
	API_KEY_IS_NIL:    "api key is nil",
}

const (
	SUCCESS            = 1
	PARA_ERROR         = 40000
	INTER_ERROR        = 40001
	SQL_ERROR          = 40002
	VERIFY_TOKEN_ERROR = 40003
	JWT_VERIFY_FAILED  = 40004
	API_KEY_IS_NIL     = 40005
)
