package core

import (
	"encoding/json"
	"errors"

	"fmt"
	"github.com/ontio/sagapi/common"
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

func (this *SagaApi) GenerateApiTestKey(apiId uint32, ontid string) (*tables.APIKey, error) {
	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		return nil, errl
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	testKey, err := dao.DefSagaApiDB.QueryApiTestKeyByOntidAndApiId(tx, ontid, apiId)
	if err != nil {
		errl = err
		return nil, err
	}

	if err != nil && !dao.IsErrNoRows(err) {
		errl = err
		return nil, err
	} else if err != nil {
		return testKey, nil
	} else {
		apiKey := &tables.APIKey{
			ApiKey:       common.GenerateUUId(common.UUID_TYPE_TEST_API_KEY),
			ApiId:        apiId,
			RequestLimit: sagaconfig.DefRequestLimit,
			UsedNum:      0,
			OntId:        ontid,
		}
		err = dao.DefSagaApiDB.InsertApiTestKey(tx, apiKey)
		if err != nil {
			errl = err
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			errl = err
			return nil, err
		}
		return apiKey, nil
	}
}

func (this *SagaApi) TestApiKey(params []tables.RequestParam) ([]byte, error) {
	if len(params) == 0 {
		return nil, errors.New("param is nil")
	}
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiId != params[i+1].ApiId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	apiTestKey := params[len(params)-1].ValueDesc
	key, err := dao.DefSagaApiDB.QueryApiKeyByApiKey(nil, apiTestKey)
	if err != nil {
		return nil, err
	}

	if key.ApiId != params[0].ApiId {
		return nil, fmt.Errorf("apiKey:%s can not invoke this api", apiTestKey)
	}

	switch key.ApiId {
	case 1:
		return this.Nasa.ApodParams(params)
	case 2:
		return this.Nasa.FeedParams(params)
	}
	return nil, fmt.Errorf("not support api, apiId:%d", key.ApiId)
}

func (this *SagaApi) QueryBasicApiInfoByPage(pageNum, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.QueryApiBasicInfoByPage(start, pageSize)
}

func (this *SagaApi) QueryBasicApiInfoByCategory(id, pageNum, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.QueryApiBasicInfoByCategoryId(nil, id, start, pageSize)
}

func (this *SagaApi) QueryApiDetailInfoByApiId(apiId uint32) (*common.ApiDetailResponse, error) {
	basicInfo, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, apiId)
	if err != nil {
		return nil, err
	}

	requestParam, err := dao.DefSagaApiDB.QueryRequestParamByApiId(nil, basicInfo.ApiId)
	if err != nil {
		return nil, err
	}

	errCode, err := dao.DefSagaApiDB.QueryErrorCode(nil)
	if err != nil {
		return nil, err
	}

	sp, err := dao.DefSagaApiDB.QuerySpecificationsByApiId(nil, basicInfo.ApiId)
	if err != nil {
		return nil, err
	}

	return &common.ApiDetailResponse{
		ApiId:               basicInfo.ApiId,
		Mark:                basicInfo.Mark,
		ResponseType:        basicInfo.RequestType,
		ResponseParam:       basicInfo.ResponseParam,
		ResponseExample:     basicInfo.ResponseExample,
		DataDesc:            basicInfo.DataDesc,
		DataSource:          basicInfo.DataSource,
		ApplicationScenario: basicInfo.ApplicationScenario,
		RequestParams:       requestParam,
		ErrorCodes:          errCode,
		Specifications:      sp,
		ApiBasicInfo:        basicInfo,
	}, nil
}

func (this *SagaApi) SearchApiIdByCategoryId(categoryId, pageNum, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.QueryApiBasicInfoByCategoryId(nil, categoryId, start, pageSize)
}

func (this *SagaApi) SearchApi() (map[string][]*tables.ApiBasicInfo, error) {
	return dao.DefSagaApiDB.SearchApi(nil)
}

type PublishErrorCode struct {
	Code int32  `json:"code"`
	Desc string `json:"description"`
}

type PublishAPI struct {
	Name            string                  `json:"name"`
	Desc            string                  `json:"description"`
	RequestType     string                  `json:"requestType"`
	ApiProvider     string                  `json:"apiProvider"`
	DataSource      string                  `json:"dataSource"`
	ResponseExample string                  `json:"responseExample"`
	Tags            []tables.Tag            `json:"tags"`
	ErrorCodes      []PublishErrorCode      `json:"errorCodes"`
	Params          []tables.RequestParam   `json:"params"`
	Specs           []tables.Specifications `json:"specifications"`
}

func PublishAPIHandleCore(param *PublishAPI) error {
	// handle error
	errorDesc, err := json.Marshal(param.ErrorCodes)
	if err != nil {
		return err
	}

	if len(param.Tags) > 100 || len(param.Params) > 100 || len(param.Specs) > 100 || len(param.ErrorCodes) > 100 {
		return err
	}

	tags := make([]*tables.Tag, 0)

	for _, tag := range param.Tags {
		t, err := dao.DefSagaApiDB.QueryTagByNameId(nil, tag.CategoryId, tag.Name)
		if err != nil {
			return err
		}
		tags = append(tags, t)
	}

	apibasic := &tables.ApiBasicInfo{
		Coin:            "ONG",
		Title:           param.Name,
		ApiProvider:     param.ApiProvider,
		ApiSagaUrlKey:   common.GenerateUUId(common.UUID_TYPE_SAGA_URL),
		ApiDesc:         param.Desc,
		ApiState:        tables.API_STATE_PUBLISH,
		ErrorDesc:       string(errorDesc),
		RequestType:     param.RequestType,
		ResponseExample: param.ResponseExample,
		DataSource:      param.DataSource,
	}

	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		return err
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	err = dao.DefSagaApiDB.InsertApiBasicInfo(tx, []*tables.ApiBasicInfo{apibasic})
	if err != nil {
		errl = err
		return err
	}

	info, err := dao.DefSagaApiDB.QueryApiBasicInfoBySagaUrlKey(tx, apibasic.ApiSagaUrlKey)
	if err != nil {
		errl = err
		return err
	}
	// tag handle

	for _, apiTag := range tags {
		tag := &tables.ApiTag{
			ApiId: info.ApiId,
			TagId: apiTag.Id,
			State: byte(1),
		}
		err = dao.DefSagaApiDB.InsertApiTag(tx, tag)
		if err != nil {
			errl = err
			//common.WriteResponse(c, common.ResponseFailed(common.SQL_ERROR, err))
			return err
		}
	}

	// handle param
	for _, p := range param.Params {
		p.ApiId = info.ApiId
		err := dao.DefSagaApiDB.InsertRequestParam(tx, []*tables.RequestParam{&p})
		if err != nil {
			errl = err
			return err
		}
	}

	// spec handle.
	for _, s := range param.Specs {
		s.ApiId = info.ApiId
		err := dao.DefSagaApiDB.InsertSpecifications(tx, []*tables.Specifications{&s})
		if err != nil {
			errl = err
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		errl = err
		return err
	}

	return nil
}
