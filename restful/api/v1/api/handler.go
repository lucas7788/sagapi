package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/restful/api/common"
	"strconv"
)

type GetBasicApiInfoByPageParam struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

func GetBasicApiInfoByPage(c *gin.Context) {
	arr, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[GetBasicApiInfoByPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageNum, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Errorf("[GetBasicApiInfoByPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Errorf("[GetBasicApiInfoByPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	infos, err := dao.DefDB.QueryApiBasicInfoByPage(start, pageSize)
	if err != nil {
		log.Errorf("[GetBasicApiInfoByPage] QueryApiBasicInfoByPage error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}

func GetApiDetailByApiId(c *gin.Context) {
	res, err := common.ParseGetParamByParamName(c, "apiId")
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	info, err := dao.DefDB.QueryApiBasicInfoByApiId(uint(apiId))
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] QueryApiBasicInfoByApiId error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(info))
}

func SearchApiByKey(c *gin.Context) {
	param, err := common.ParseGetParamByParamName(c, "key")
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	infos, err := dao.DefDB.SearchApi(param[0])
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] SearchApi error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}
