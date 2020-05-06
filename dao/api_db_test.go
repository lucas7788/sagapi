package dao

import (
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApiDB_InsertApiBasicInfo(t *testing.T) {
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
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
	l := 11
	infos := make([]*tables.ApiBasicInfo, l)
	for i := 0; i < len(infos); i++ {
		infos[i] = info
	}
	assert.Nil(t, TestDB.ApiDB.InsertApiBasicInfo(infos))

	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(1, 1)
	assert.Nil(t, err)
	info, err = TestDB.ApiDB.QueryApiBasicInfoByApiId(basic[0].ApiId)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, basic[0].ApiId)

	detail := &tables.ApiDetailInfo{
		ApiId:               basic[0].ApiId,
		RequestType:         "POST",
		Mark:                "",
		ResponseParam:       "",
		ResponseExample:     "",
		DataDesc:            "",
		DataSource:          "",
		ApplicationScenario: "",
	}
	err = TestDB.ApiDB.InsertApiDetailInfo(detail)
	assert.Nil(t, err)
	detail, err = TestDB.ApiDB.QueryApiDetailInfoByApiId(basic[0].ApiId)
	assert.Nil(t, err)
	params := &tables.RequestParam{
		ApiDetailInfoId: detail.Id,
		ParamName:       "",
		Required:        true,
		ParamType:       "",
		Note:            "",
		ValueDesc:       "",
	}

	assert.Nil(t, TestDB.ApiDB.InsertRequestParam([]*tables.RequestParam{params}))

	requestParams, err := TestDB.ApiDB.QueryRequestParamByApiDetailId(detail.Id)
	assert.Nil(t, err)
	assert.Equal(t, len(requestParams), 1)

	infos, err = TestDB.ApiDB.QueryApiBasicInfoByCategoryId(1, 0, 1)
	assert.Nil(t, err)

	TestDB.ApiDB.ClearRequestParamDB()
	TestDB.ApiDB.ClearApiBasicDB()
}

func TestApiDB_InsertApiKey(t *testing.T) {

	TestDB.ApiDB.ClearApiKeyDB()
	TestDB.OrderDB.ClearOrderDB()
	TestDB.ApiDB.ClearApiBasicDB()

	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)

	orderId := "orderId"
	tt := time.Now().Unix()
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.OrderDB.InsertOrder(order)
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

	err = TestDB.ApiDB.InsertApiKey(key)
	assert.Nil(t, err)

	key, err = TestDB.ApiDB.QueryApiKeyByApiKey(apikey)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), key.UsedNum)

	TestDB.ApiDB.ClearApiKeyDB()
	TestDB.OrderDB.ClearOrderDB()
	TestDB.ApiDB.ClearApiBasicDB()
}

func TestApiDB_QuerySpecificationsByApiDetailId(t *testing.T) {
	assert.Nil(t, TestDB.ApiDB.ClearSpecificationsDB())
	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)
	detail, err := TestDB.ApiDB.QueryApiDetailInfoByApiId(basic[0].ApiId)
	assert.Nil(t, err)
	params := []*tables.Specifications{
		&tables.Specifications{
			ApiDetailInfoId: detail.Id,
			Price:           "0",
			Amount:          500,
		},
		&tables.Specifications{
			ApiDetailInfoId: detail.Id,
			Price:           "0.01",
			Amount:          1000,
		},
	}
	err = TestDB.ApiDB.InsertSpecifications(params)
	assert.Nil(t, err)

	specs, err := TestDB.ApiDB.QuerySpecificationsByApiDetailId(detail.Id)
	assert.Nil(t, err)
	assert.Equal(t, specs[0].ApiDetailInfoId, detail.Id)

	spec, err := TestDB.ApiDB.QuerySpecificationsById(specs[0].Id)
	assert.Nil(t, err)

	assert.Equal(t, spec.Id, specs[0].Id)

	assert.Nil(t, TestDB.ApiDB.ClearSpecificationsDB())
}
