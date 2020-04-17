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
		ApiId:1,
		OrderId:   "abc",
		OntId:     "111",
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
	key, err := TestDB.ApiDB.QueryApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, 1, key.UsedNum)
}

func TestSagaDB_QueryRequestNum(t *testing.T) {
	key, err := TestDB.ApiDB.QueryApiKey("apikey")
	assert.Nil(t, err)
	assert.Equal(t, 1, key.UsedNum)
}

func TestSagaDB_SearchApi(t *testing.T) {
	info := &tables.ApiBasicInfo{
		ApiDesc:        "abcdefg",
		ApiPrice:       "0.1",
		Specifications: 1,
	}
	info2 := &tables.ApiBasicInfo{
		ApiDesc:        "cdefgty",
		ApiPrice:       "0.1",
		Specifications: 1,
	}
	err := TestDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info})
	assert.Nil(t, err)
	err = TestDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info2})
	assert.Nil(t, err)
	infos, err := TestDB.ApiDB.SearchApi("cdefgty")
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
