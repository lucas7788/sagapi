package dao

import (
	"testing"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"time"
)

var TestDB *SagaApiDB

func TestMain(m *testing.M) {
	fmt.Println("begin test db.")
	var err error
	sagaDBConfig := config.DefSagaConfig
	sagaDBConfig.DbConfig = config.DefDBConfigMap[config.NETWORK_ID_TRAVIS_NET]
	TestDB, err = NewSagaApiDB(sagaDBConfig)
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Println("end")
}

func TestOrderDB_UpdateOrderStatus(t *testing.T) {
	err := TestDB.OrderDB.UpdateOrderStatus("cc5ebaf6-24cb-423b-a590-e05a41f8c1f5", config.Canceled)
	assert.Nil(t, err)
}

func TestOrderDB_InsertOrder(t *testing.T) {
	tt := time.Now().Unix()
	orderId := "abcedkfy"
	br := &tables.Order{
		ApiId:     1,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err := TestDB.OrderDB.InsertOrder(br)
	assert.Nil(t, err)
	code := &tables.QrCode{
		QrCodeId: "qbcdab",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.OrderDB.InsertQrCode(code)
	assert.Nil(t, err)
	code = &tables.QrCode{
		QrCodeId: "qbcdabc",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.OrderDB.InsertQrCode(code)
	assert.Nil(t, err)
	code, err = TestDB.OrderDB.QueryQrCodeByOrderId(orderId)
	assert.Nil(t, err)
	fmt.Println(code)
	err = TestDB.OrderDB.DeleteQrCodeByOrderId(orderId)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId)
	assert.Nil(t, err)
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
}

func TestSagaDB_QueryRequestNum(t *testing.T) {
	key, err := TestDB.ApiDB.QueryApiKeyByApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, 1, key.UsedNum)
}

func TestSagaDB_SearchApi(t *testing.T) {
	info := &tables.ApiBasicInfo{
		ApiDesc:        "abcdefg",
		Price:          "0.1",
		Specifications: 1,
	}
	info2 := &tables.ApiBasicInfo{
		ApiDesc:        "cdefgty",
		Price:          "0.1",
		Specifications: 1,
	}
	err := TestDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info})
	assert.Nil(t, err)
	err = TestDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info2})
	assert.Nil(t, err)
	infos, err := TestDB.ApiDB.SearchApiByKey("cdefgty")
	assert.Nil(t, err)
	fmt.Println(infos)
	infos, err = TestDB.ApiDB.QueryApiBasicInfoByPage(2, 2)
	assert.Nil(t, err)
	fmt.Println(infos)
	info3, err := TestDB.ApiDB.QueryApiBasicInfoByApiId(100)
	assert.NotNil(t, err)
	fmt.Println(info3)
}

func TestSagaDB_QueryOrderStatusByOrderId(t *testing.T) {
	_, err := TestDB.OrderDB.QueryOrderByOrderId("1")
	assert.NotNil(t, err)
}

func TestOrderDB_QueryOrderByPage(t *testing.T) {
	order, err := TestDB.OrderDB.QueryOrderByPage(1, 2, "")
	assert.Nil(t, err)
	fmt.Println(order)
}

func TestOrderDB_UpdateTxInfoByOrderId(t *testing.T) {
	err := TestDB.OrderDB.UpdateTxInfoByOrderId("145e89f6-850e-44a7-be3e-9224fd066858", "", config.Processing)
	assert.Nil(t, err)
}

func TestOrderDB_QueryQrCodeResultByQrCodeId(t *testing.T) {
	status, err := TestDB.OrderDB.QueryQrCodeResultByQrCodeId("bb5b68d4-0282-469d-936b-ae43e30c5de5")
	assert.Nil(t, err)
	fmt.Println(status)
}
