package core

import (
	"errors"

	"fmt"
	"github.com/ontio/sagapi/common"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/core/nasa"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

var DefSagaApi *SagaApi

type SagaApi struct {
	Nasa      *nasa.Nasa
	SagaOrder *SagaOrder
}

func NewSagaApi() *SagaApi {
	http.DefClient = http.NewClient()
	return &SagaApi{
		Nasa:      nasa.NewNasa(),
		SagaOrder: NewSagaOrder(),
	}
}

func (this *SagaApi) GenerateApiTestKey(apiId int, ontid string) (*tables.APIKey, error) {
	var testKey tables.APIKey
	err := dao.DefApiDb.Conn.Get(testKey, "select * from tbl_api_test_key where OntId=? and ApiId=?", ontid, apiId)
	if err != nil {
		return nil, err
	}

	if err != nil && !dao.IsNoEltError(err) {
		return nil, err
	}

	key := "test_" + common.GenerateUUId()
	apiKey := &tables.APIKey{
		ApiKey:       key,
		ApiId:        apiId,
		RequestLimit: sagaconfig.DefRequestLimit,
		UsedNum:      0,
		OntId:        ontid,
	}
	err = dao.DefSagaApiDB.ApiDB.InsertApiTestKey(apiKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (this *SagaApi) TestApiKey(params []tables.RequestParam) ([]byte, error) {
	if len(params) == 0 {
		return nil, errors.New("param is nil")
	}
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiDetailInfoId != params[i+1].ApiDetailInfoId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	apiTestKey := params[len(params)-1].ValueDesc

	var key tables.APIKey
	err := dao.DefApiDb.Conn.Get(&key, "select * from tbl_api_test_key where where ApiKey=?", apiTestKey)
	if err != nil {
		return nil, err
	}

	switch key.ApiId {
	case 1:
		return this.Nasa.ApodParams(params)
	case 2:
		return this.Nasa.FeedParams(params)
	}
	return nil, fmt.Errorf("not support api, apiId:%d", key.ApiId)
}

func (this *SagaApi) QueryBasicApiInfoByPage(pageNum, pageSize int) ([]tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize

	var result []tables.ApiBasicInfo
	err := dao.DefApiDb.Conn.Select(&result, "select * from tbl_api_basic_info where ApiId limit ?, ?", start, pageSize)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *SagaApi) QueryBasicApiInfoByCategory(id, pageNum, pageSize int) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoByCategoryId(id, start, pageSize)
}

func (this *SagaApi) QueryApiDetailInfoByApiId(apiId int) (*common.ApiDetailResponse, error) {
	var basicInfo tables.ApiBasicInfo
	err := dao.DefApiDb.Conn.Get(&basicInfo, "select * from tbl_api_basic_info where ApiId=?", apiId)
	if err != nil {
		return nil, err
	}

	var apiDetail tables.ApiDetailInfo
	err = dao.DefApiDb.Conn.Get(&apiDetail, "select * from tbl_api_detail_info where ApiId=?", apiId)
	if err != nil {
		return nil, err
	}

	var requestParam []tables.RequestParam
	err = dao.DefApiDb.Conn.Select(&requestParam, "select * from tbl_request_param where ApiDetailInfoId=?", apiDetail.Id)
	if err != nil {
		return nil, err
	}

	var errCode []tables.ErrorCode
	err = dao.DefApiDb.Conn.Select(&errCode, "select * from tbl_error_code where ApiDetailInfoId=?", apiDetail.Id)
	if err != nil {
		return nil, err
	}

	var sp []tables.Specifications
	err = dao.DefApiDb.Conn.Select(&sp, "select * from tbl_specifications where ApiDetailInfoId=?", apiDetail.Id)
	if err != nil {
		return nil, err
	}

	return &common.ApiDetailResponse{
		ApiId:               apiDetail.ApiId,
		Mark:                apiDetail.Mark,
		ResponseType:        apiDetail.RequestType,
		ResponseParam:       apiDetail.ResponseParam,
		ResponseExample:     apiDetail.ResponseExample,
		DataDesc:            apiDetail.DataDesc,
		DataSource:          apiDetail.DataSource,
		ApplicationScenario: apiDetail.ApplicationScenario,
		RequestParams:       requestParam,
		ErrorCodes:          errCode,
		Specifications:      sp,
		ApiBasicInfo:        &basicInfo,
	}, nil
}

func (this *SagaApi) SearchApiIdByCategoryId(categoryId, pageNum, pageSize int) ([]tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	var result []tables.ApiBasicInfo
	if categoryId == 1 {
		err := dao.DefApiDb.Conn.Select(&result, "select * from tbl_api_basic_info where ApiId limit ?, ?", start, pageSize)
		if err != nil {
			return nil, err
		}
	} else {
		err := dao.DefApiDb.Conn.Select(&result, "select * from tbl_api_basic_info where ApiId in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where category_id=?))", categoryId)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *SagaApi) SearchApi() (map[string][]tables.ApiBasicInfo, error) {
	res := make(map[string][]tables.ApiBasicInfo)
	var newestApi []tables.ApiBasicInfo
	var hottestApi []tables.ApiBasicInfo
	var freeApi []tables.ApiBasicInfo

	err := dao.DefApiDb.Conn.Select(&newestApi, "select * from tbl_api_basic_info order by CreateTime limit ?", 10)
	if err != nil {
		return nil, err
	}
	res["newest"] = newestApi

	err = dao.DefApiDb.Conn.Select(&hottestApi, "select * from tbl_api_basic_info order by InvokeFrequency limit ?", 10)
	if err != nil {
		return nil, err
	}
	res["hottest"] = hottestApi

	err = dao.DefApiDb.Conn.Select(&freeApi, "select * from tbl_api_basic_info where Price='0' limit ?", 10)
	if err != nil {
		return nil, err
	}
	res["free"] = freeApi
	return res, nil
}
