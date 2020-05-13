package core

import (
	"encoding/json"
	"errors"

	"fmt"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core/freq"
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
	Cache     *freq.DBCache
}

func NewSagaApi() *SagaApi {
	http.DefClient = http.NewClient()
	cache := freq.NewDBCache()
	return &SagaApi{
		Nasa:      nasa.NewNasa(cache),
		SagaOrder: NewSagaOrder(),
		Cache:     cache,
	}
}

func (this *SagaApi) GenerateApiTestKey(apiId uint32, ontid string, apiState int32) (*tables.APIKey, error) {
	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("GenerateApiTestKey.0. %s", errl)
		return nil, errl
	}

	_, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(tx, apiId, apiState)
	if err != nil {
		errl = err
		log.Debugf("GenerateApiTestKey.1. %s", err)
		return nil, err
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	testKey, err := dao.DefSagaApiDB.QueryApiTestKeyByOntidAndApiId(tx, ontid, apiId)
	if err != nil && !dao.IsErrNoRows(err) {
		errl = err
		log.Debugf("GenerateApiTestKey.3. %s", err)
		return nil, err
	} else if err == nil {
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

func (this *SagaApi) AdminTestApi(params []*tables.RequestParam, apiId uint32) ([]byte, error) {
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiId != params[i+1].ApiId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].Required && params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	info, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, apiId, tables.API_STATE_PUBLISH)
	if err != nil {
		return nil, err
	}

	data, err := HandleDataSourceReqCore(nil, info.ApiSagaUrlKey, params, "", true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *SagaApi) TestApiKey(params []*tables.RequestParam, apiKey string) ([]byte, error) {
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiId != params[i+1].ApiId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].Required && params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	apiTestKey := apiKey

	key, err := dao.DefSagaApiDB.QueryApiKeyByApiKey(nil, apiTestKey)
	if err != nil {
		return nil, err
	}

	apiId := key.ApiId
	switch apiId {
	case 1:
		if len(params) != 0 {
			return nil, fmt.Errorf("apod no need params")
		}
		return this.Nasa.ApodParams(apiTestKey)
	case 2:
		return this.Nasa.FeedParams(params, apiTestKey)
	default:
		info, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, apiId, tables.API_STATE_BUILTIN)
		if err != nil {
			return nil, err
		}

		data, err := HandleDataSourceReqCore(nil, info.ApiSagaUrlKey, params, apiTestKey, false)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (this *SagaApi) QueryBasicApiInfoByPage(pageNum, pageSize uint32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefSagaApiDB.QueryApiBasicInfoByPage(start, pageSize, apiState)
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

func (this *SagaApi) QueryApiDetailInfoByApiId(apiId uint32, apiState int32) (*common.ApiDetailResponse, error) {
	basicInfo, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, apiId, apiState)
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
	Coin            string                  `json:"coin"`
	Tags            []tables.Tag            `json:"tags"`
	ErrorCodes      []PublishErrorCode      `json:"errorCodes"`
	Params          []tables.RequestParam   `json:"params"`
	Specs           []tables.Specifications `json:"specifications"`
}

func PublishAPIHandleCore(param *PublishAPI, ontId, author string) error {
	// handle error
	errorDesc, err := json.Marshal(param.ErrorCodes)
	if err != nil {
		log.Debugf("PublishAPIHandleCore.0. %s", err)
		return err
	}

	if len(param.Tags) > 100 || len(param.Params) > 100 || len(param.Specs) > 100 || len(param.ErrorCodes) > 100 {
		log.Debugf("PublishAPIHandleCore.1. %s", err)
		return err
	}

	if len(param.Tags) == 0 {
		log.Debugf("PublishAPIHandleCore.1. tag can not empty")
		return errors.New("PublishAPIHandleCore.1. tag can not empty")
	}

	tags := make([]*tables.Tag, 0)

	for _, tag := range param.Tags {
		t, err := dao.DefSagaApiDB.QueryTagByNameId(nil, tag.CategoryId, tag.Name)
		if err != nil {
			log.Debugf("PublishAPIHandleCore.2. %s", err)
			return err
		}
		tags = append(tags, t)
	}

	cat, err := dao.DefSagaApiDB.QueryCategoryById(nil, tags[0].CategoryId)
	if err != nil {
		log.Debugf("PublishAPIHandleCore.2.0 %s", err)
		return err
	}

	if param.Coin != "ONG" && param.Coin != "ONT" {
		return errors.New("wrong Coin type. only ONT/ONG")
	}

	if param.RequestType != "POST" && param.RequestType != "GET" {
		return errors.New("wrong RequestType type. only POST/GET")
	}

	apibasic := &tables.ApiBasicInfo{
		Coin:            param.Coin,
		Title:           param.Name,
		Icon:            cat.Icon,
		ApiProvider:     param.ApiProvider,
		ApiSagaUrlKey:   common.GenerateUUId(common.UUID_TYPE_SAGA_URL),
		ApiDesc:         param.Desc,
		ApiState:        tables.API_STATE_PUBLISH,
		ErrorDesc:       string(errorDesc),
		RequestType:     param.RequestType,
		ResponseExample: param.ResponseExample,
		DataSource:      param.DataSource,
		OntId:           ontId,
		Author:          author,
	}
	port := fmt.Sprintf("%d", sagaconfig.DefSagaConfig.RestPort)
	apibasic.ApiUrl = sagaconfig.DefSagaConfig.SagaHost + ":" + port + "/api/v1/data_source/" + apibasic.ApiSagaUrlKey + "/:apikey"

	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("PublishAPIHandleCore.3. %s", err)
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
		log.Debugf("PublishAPIHandleCore.4. %s", err)
		return err
	}

	info, err := dao.DefSagaApiDB.QueryApiBasicInfoBySagaUrlKey(tx, apibasic.ApiSagaUrlKey, tables.API_STATE_PUBLISH)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.5. %s", err)
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
			log.Debugf("PublishAPIHandleCore.6. %s", err)
			return err
		}
	}

	// handle param
	for _, p := range param.Params {
		p.ApiId = info.ApiId
		err := dao.DefSagaApiDB.InsertRequestParam(tx, []*tables.RequestParam{&p})
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.7. %s", err)
			return err
		}
	}

	// spec handle.
	for _, s := range param.Specs {
		s.ApiId = info.ApiId
		err := dao.DefSagaApiDB.InsertSpecifications(tx, []*tables.Specifications{&s})
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.8. %s", err)
			return err
		}
	}

	referParams, err := dao.DefSagaApiDB.QueryRequestParamByApiId(tx, info.ApiId)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.9. %s", err)
		return err
	}
	_, err = HandleDataSourceReqCore(tx, info.ApiSagaUrlKey, referParams, "", true)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.10. %s", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Debugf("PublishAPIHandleCore.11. %s", err)
		errl = err
		return err
	}

	return nil
}
