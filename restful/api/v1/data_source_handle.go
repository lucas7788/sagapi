package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
)

func HandleDataSourceReq(c *gin.Context) {
	var params []*tables.RequestParam
	sagaUrlKey := c.Param("sagaUrlKey")
	if sagaUrlKey == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("url key can not empty")))
		return
	}
	apiKey := c.Param("apikey")
	if apiKey == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("apikey is nil")))
		return
	}

	err := common.ParsePostParam(c, &params)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
	}

	data, err := core.HandleDataSourceReqCore(nil, sagaUrlKey, params, apiKey, false)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(string(data)))
}
