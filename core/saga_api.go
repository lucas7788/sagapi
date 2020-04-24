package core

import (
	"errors"

	"fmt"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/core/nasa"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"strings"
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
	testKey, err := this.QueryApiTestKeyByOntIdAndApiId(ontid, apiId)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}
	if testKey != nil {
		return testKey, nil
	}
	key := "test_" + common.GenerateUUId()
	apiKey := &tables.APIKey{
		ApiKey:       key,
		ApiId:        apiId,
		RequestLimit: 10,
		UsedNum:      0,
		OntId:        ontid,
	}
	err = dao.DefSagaApiDB.ApiDB.InsertApiTestKey(apiKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (this *SagaApi) TestApiKey(params []tables.RequestParam) (string, error) {
	if len(params) == 0 {
		return "", errors.New("param is nil")
	}
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiDetailInfoId != params[i+1].ApiDetailInfoId {
			return "", errors.New("params should to the same api")
		}
		if params[i].ValueDesc == "" {
			return "", fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	apiTestKey := params[len(params)-1].ValueDesc
	key, err := dao.DefSagaApiDB.ApiDB.QueryApiTestKeyByApiTestKey(apiTestKey)
	if err != nil {
		return "", err
	}

	switch key.ApiId {
	case 1:
		return this.Nasa.ApodParams(params)
	case 2:
		return this.Nasa.FeedParams(params)
	}
	return "", fmt.Errorf("not support api, apiId:%d", key.ApiId)
}

func (this *SagaApi) QueryApiTestKeyByOntIdAndApiId(ontid string, apiId int) (*tables.APIKey, error) {
	return dao.DefSagaApiDB.ApiDB.QueryApiTestKeyByOntIdAndApiId(ontid, apiId)
}

func (this *SagaApi) QueryBasicApiInfoByPage(pageNum, pageSize int) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoByPage(start, pageSize)
}

func (this *SagaApi) QueryApiDetailInfoByApiId(apiId int) (*common.ApiDetailResponse, error) {
	basicInfo, err := dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoByApiId(apiId)
	if err != nil {
		return nil, err
	}
	apiDetail, err := dao.DefSagaApiDB.ApiDB.QueryApiDetailInfoById(apiId)
	if err != nil {
		return nil, err
	}
	if apiDetail == nil {
		return nil, nil
	}
	requestParam, err := dao.DefSagaApiDB.ApiDB.QueryRequestParamByApiDetailInfoId(apiDetail.Id)
	if err != nil {
		return nil, err
	}
	errCode, err := dao.DefSagaApiDB.ApiDB.QueryErrorCodeByApiDetailInfoId(apiDetail.Id)
	if err != nil {
		return nil, err
	}
	sp, err := dao.DefSagaApiDB.ApiDB.QuerySpecificationsByApiDetailId(apiDetail.Id)
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
		ApiBasicInfo:        basicInfo,
	}, nil
}

func (this *SagaApi) SearchApiIdByCategoryId(categoryId int) ([]*tables.ApiBasicInfo, error) {

	return dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoByCategoryId(categoryId)
}

//newest hot free
func (this *SagaApi) SearchApi() (map[string][]*tables.ApiBasicInfo, error) {
	res := make(map[string][]*tables.ApiBasicInfo)
	//newest
	newest, err := dao.DefSagaApiDB.ApiDB.QueryNewestApiBasicInfo()
	if err != nil {
		return nil, err
	}
	res["newest"] = newest
	hottest, err := dao.DefSagaApiDB.ApiDB.QueryHottestApiBasicInfo()
	if err != nil {
		return nil, err
	}
	res["hottest"] = hottest
	free, err := dao.DefSagaApiDB.ApiDB.QueryFreeApiBasicInfo()
	res["free"] = free
	return res, nil
}
