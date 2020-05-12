package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
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
	infos, err := core.DefSagaApi.QueryBasicApiInfoByPage(uint32(pageNum), uint32(pageSize), tables.API_STATE_BUILTIN)
	if err != nil {
		log.Errorf("[GetBasicApiInfoByPage] QueryApiBasicInfoByPage error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	res := common.PageResult{
		List:  infos,
		Count: uint32(len(infos)),
	}

	common.WriteResponse(c, common.ResponseSuccess(res))
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
	info, err := core.DefSagaApi.QueryApiDetailInfoByApiId(uint32(apiId), tables.API_STATE_BUILTIN)
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] QueryApiDetailInfoByApiId error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(info))
}

func SearchApiByKey(c *gin.Context) {
	key := &common2.SearchApiByKey{}
	err := common.ParsePostParam(c, key)
	if err != nil {
		log.Errorf("[SearchApiByKey] ParsePostParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if key == nil || key.Key == "" {
		common.WriteResponse(c, common.ResponseSuccess(nil))
		return
	}
	//todo key.Key should not have sql statement
	var infonil []*tables.ApiBasicInfo
	infos, err := dao.DefSagaApiDB.SearchApiByKey(key.Key)
	if err != nil && dao.IsErrNoRows(err) {
		// empty
		common.WriteResponse(c, common.ResponseSuccess(infonil))
		return
	} else if err != nil {
		log.Errorf("[GetApiDetailByApiId] SearchApiByKey error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}

func SearchApiByCategoryId(c *gin.Context) {
	param := &common2.GetApiByCategoryId{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[SearchApiByCategoryId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	infos, err := core.DefSagaApi.SearchApiIdByCategoryId(param.CategoryId, param.PageNum, param.PageSize)
	if err != nil {
		log.Errorf("[SearchApiByCategoryId] SearchApiByKey error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	log.Debugf("SearchApiByCategoryId: num %d", len(infos))
	common.WriteResponse(c, common.ResponseSuccess(infos))
}

func SearchApi(c *gin.Context) {
	infos, err := core.DefSagaApi.SearchApi()
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] SearchApiByKey error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}
