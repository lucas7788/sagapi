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

	infos, err := TestDB.ApiDB.QueryApiBasicInfoByCategoryId(10, 0,1)
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

	kfre,err := TestDB.ApiDB.QueryApiKeyAndInvokeFreByApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, kfre.ApiKey, "apikey")
}

func TestSagaDB_QueryRequestNum(t *testing.T) {
	key, err := TestDB.ApiDB.QueryApiKeyByApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, 1, key.UsedNum)
}