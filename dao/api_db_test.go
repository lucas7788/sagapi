package dao

import (
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApiDB_InsertApiBasicInfo(t *testing.T) {
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
		ApiProvider:     common.GenerateUUId(),
		ApiUrl:          "",
		Price:           "",
		ApiDesc:         "",
		Specifications:  1,
		Popularity:      0,
		Delay:           0,
		SuccessRate:     0,
		InvokeFrequency: 0,
	}
	l := 11
	infos := make([]*tables.ApiBasicInfo, l)
	for i := 0; i < len(infos); i++ {
		info.ApiProvider = common.GenerateUUId()
		infos[i] = info
	}
	assert.Nil(t, TestDB.InsertApiBasicInfo(infos))

	Coin := "ONG"
	ApiType := "ApiType"
	Icon := "Icon"
	Title := "Title"
	ApiProvider := "ApiProvider"
	ApiSagaUrlKey := "ApiSagaUrlKey"
	ApiUrl := "ApiUrl"
	Price := "Price"
	ApiDesc := "ApiDesc"
	ErrorDesc := "ErrorDesc"
	Specifications := uint32(9)
	Popularity := uint32(10)
	Delay := uint32(11)
	SuccessRate := uint32(100)
	InvokeFrequency := uint64(12345)
	ApiState := int32(99)
	RequestType := "RequestType"
	Mark := "Mark"
	ResponseParam := "ResponseParam"
	ResponseExample := "ResponseExample"
	DataDesc := "DataDesc"
	DataSource := "DataSource"
	ApplicationScenario := "ApplicationScenario"

	info2 := &tables.ApiBasicInfo{
		Coin:                Coin,
		ApiType:             ApiType,
		Icon:                Icon,
		Title:               Title,
		ApiProvider:         ApiProvider,
		ApiSagaUrlKey:       ApiSagaUrlKey,
		ApiUrl:              ApiUrl,
		Price:               Price,
		ApiDesc:             ApiDesc,
		ErrorDesc:           ErrorDesc,
		Specifications:      Specifications,
		Popularity:          Popularity,
		Delay:               Delay,
		SuccessRate:         SuccessRate,
		InvokeFrequency:     InvokeFrequency,
		ApiState:            ApiState,
		RequestType:         RequestType,
		Mark:                Mark,
		ResponseParam:       ResponseParam,
		ResponseExample:     ResponseExample,
		DataDesc:            DataDesc,
		DataSource:          DataSource,
		ApplicationScenario: ApplicationScenario,
	}

	assert.Nil(t, TestDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info2}))

	basic, err := TestDB.QueryApiBasicInfoByPage(1, 1)
	assert.Nil(t, err)
	info, err = TestDB.QueryApiBasicInfoByApiId(basic[0].ApiId)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, basic[0].ApiId)

	assert.Nil(t, err)
	params := &tables.RequestParam{
		ApiId:      info.ApiId,
		ParamName:  "",
		Required:   true,
		ParamType:  "",
		ParamWhere: tables.URL_PARAM_RESTFUL,
		Note:       "",
		ValueDesc:  "",
	}

	assert.Nil(t, TestDB.InsertRequestParam([]*tables.RequestParam{params}))

	requestParams, err := TestDB.QueryRequestParamByApiDetailId(detail.Id)
	assert.Nil(t, err)
	assert.Equal(t, len(requestParams), 1)

	infos, err = TestDB.QueryApiBasicInfoByCategoryId(1, 0, 1)
	assert.Nil(t, err)

	err = TestDB.ClearRequestParamDB()
	assert.Nil(t, err)
	err = TestDB.ClearApiBasicDB()
	assert.Nil(t, err)
}

func TestApiDB_InsertApiKey(t *testing.T) {

	TestDB.ClearApiKeyDB()
	TestDB.ClearOrderDB()
	TestDB.ClearApiBasicDB()

	basic, err := TestDB.QueryApiBasicInfoByPage(0, 1)

	orderId := "orderId"
	tt := time.Now().Unix()
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(order)
	assert.Nil(t, err)
	apikey := "apikey"
	key := &tables.APIKey{
		ApiId:        basic[0].ApiId,
		OrderId:      orderId,
		ApiKey:       apikey,
		RequestLimit: 2,
		UsedNum:      1,
		OntId:        "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
	}

	err = TestDB.InsertApiKey(key)
	assert.Nil(t, err)

	key, err = TestDB.QueryApiKeyByApiKey(apikey)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), key.UsedNum)

	TestDB.ClearApiKeyDB()
	TestDB.ClearOrderDB()
	TestDB.ClearApiBasicDB()
}

func TestApiDB_QuerySpecificationsByApiDetailId(t *testing.T) {
	assert.Nil(t, TestDB.ClearSpecificationsDB())
	basic, err := TestDB.QueryApiBasicInfoByPage(0, 1)
	detail, err := TestDB.QueryApiDetailInfoByApiId(basic[0].ApiId)
	assert.Nil(t, err)
	params := []*tables.Specifications{
		&tables.Specifications{
			ApiId:  detail.Id,
			Price:  "0",
			Amount: 500,
		},
		&tables.Specifications{
			ApiId:  detail.Id,
			Price:  "0.01",
			Amount: 1000,
		},
	}
	err = TestDB.InsertSpecifications(params)
	assert.Nil(t, err)

	specs, err := TestDB.QuerySpecificationsByApiDetailId(detail.Id)
	assert.Nil(t, err)
	assert.Equal(t, specs[0].ApiId, detail.Id)

	spec, err := TestDB.QuerySpecificationsById(specs[0].Id)
	assert.Nil(t, err)

	assert.Equal(t, spec.Id, specs[0].Id)

	assert.Nil(t, TestDB.ClearSpecificationsDB())
}

func TestApiDB_QueryApiBasicInfoBySagaUrlKey(t *testing.T) {
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
		ApiSagaUrlKey:   "sagaurl_" + common.GenerateUUId(),
		ApiProvider:     "",
		ApiUrl:          "",
		Price:           "",
		ApiDesc:         "",
		Specifications:  1,
		Popularity:      0,
		Delay:           0,
		SuccessRate:     0,
		InvokeFrequency: 0,
	}

	assert.Nil(t, TestDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info}))
	_, err := TestDB.QueryApiBasicInfoBySagaUrlKey(info.ApiSagaUrlKey)
	assert.Nil(t, err)
}
