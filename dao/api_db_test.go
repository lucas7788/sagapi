package dao

import (
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
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
}

func TestApiDB_QueryApiBasicInfoByApiId(t *testing.T) {
	info, err := TestDB.ApiDB.QueryApiBasicInfoByApiId(1)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, 1)

	infos, err := TestDB.ApiDB.QueryApiBasicInfoByCategoryId(10, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, len(infos), 1)
}

func TestApiDB_InsertApiKey(t *testing.T) {
	orderId := "145e89f6-850e-44a7-be3e-9224fd066858"
	key := &tables.APIKey{
		ApiId:        1,
		OrderId:      orderId,
		ApiKey:       "apikey",
		RequestLimit: 2,
		UsedNum:      1,
		OntId:        "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
	}

	err := TestDB.ApiDB.InsertApiKey(key)
	assert.NotNil(t, err)

	ord := &tables.Order{
		ApiId:   1,
		OrderId: orderId,
		OntId:   "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
	}
	err = TestDB.OrderDB.InsertOrder(ord)
	assert.Nil(t, err)
	err = TestDB.ApiDB.InsertApiKey(key)
	assert.Nil(t, err)
}

func TestApiDB_QueryApiKey(t *testing.T) {
	key, err := TestDB.ApiDB.QueryApiKeyByApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, 1, key.UsedNum)

	kfre, err := TestDB.ApiDB.QueryApiKeyAndInvokeFreByApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, "apikey", kfre.ApiKey)
}

func TestApiDB_InsertSpecifications(t *testing.T) {
	params := []*tables.Specifications{
		&tables.Specifications{
			ApiDetailInfoId: 1,
			Price:           "0",
			Amount:          500,
		},
		&tables.Specifications{
			ApiDetailInfoId: 1,
			Price:           "0.01",
			Amount:          1000,
		},
	}
	err := TestDB.ApiDB.InsertSpecifications(params)
	assert.Nil(t, err)
}

func TestApiDB_QuerySpecificationsByApiDetailId(t *testing.T) {
	spec, err := TestDB.ApiDB.QuerySpecificationsById(1)
	assert.Nil(t, err)

	assert.Equal(t, spec.Id, 1)

	specs, err := TestDB.ApiDB.QuerySpecificationsByApiDetailId(1)
	assert.Nil(t, err)
	assert.Equal(t, len(specs), 5)
}
