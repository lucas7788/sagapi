package dao

import (
	"testing"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"time"
)

var TestDB *SagaApiDB

func TestMain(m *testing.M) {
	fmt.Println("begin")
	var err error
	TestDB, err = NewSagaApiDB()
	if err != nil {
		return
	}
	m.Run()
	fmt.Println("end")
}

func TestOrderDB_InsertOrder(t *testing.T) {
	tt := time.Now().Unix()
	br := &tables.Order{
		ApiId:     1,
		OrderId:   "abcdefg",
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err := TestDB.OrderDB.InsertOrder(br)
	assert.Nil(t, err)
}

func TestApiDB_InsertApiKey(t *testing.T) {
	key := &tables.APIKey{
		ApiId:        1,
		OrderId:      "abc",
		ApiKey:       "apikey",
		RequestLimit: 2,
		UsedNum:      1,
	}
	err := TestDB.ApiDB.InsertApiKey(key)
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
	assert.Nil(t, err)
	fmt.Println(info3)
	TestDB.Close()
}

func TestSagaDB_QueryOrderStatusByOrderId(t *testing.T) {
	status, err := TestDB.OrderDB.QueryOrderByOrderId("1")
	assert.Nil(t, err)
	fmt.Println("status:", status)
}

func TestOrderDB_QueryOrderByPage(t *testing.T) {
	order, err := TestDB.OrderDB.QueryOrderByPage(1, 2, "")
	assert.Nil(t, err)
	fmt.Println(order)
}

func TestOrderDB_QueryQrCodeResultByQrCodeId(t *testing.T) {
	status, err := TestDB.OrderDB.QueryQrCodeResultByQrCodeId("bb5b68d4-0282-469d-936b-ae43e30c5de5")
	assert.Nil(t, err)
	fmt.Println(status)
}
