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

func Init(t *testing.T) {
	config.DefConfig.DbConfig.ProjectDBUrl = "127.0.0.1:3306"
	config.DefConfig.DbConfig.ProjectDBName = "saga"
	config.DefConfig.DbConfig.ProjectDBUser = "root"
	config.DefConfig.DbConfig.ProjectDBPassword = "111111"

	db, err := NewSagaApiDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.Nil(t, err)
	TestDB = db
}

func TestQueryTestRecord(t *testing.T) {
	type Email struct {
		ID         int
		UserID     int    `gorm:"index"`                          // 外键 (属于), tag `index`是为该列创建索引
		Email      string `gorm:"type:varchar(100);unique_index"` // `type`设置sql类型, `unique_index` 为该列设置唯一索引
		Subscribed bool
	}
	type User struct {
		Name   string
		Emails []Email
	}

	Init(t)
}

func TestSagaDB_Init(t *testing.T) {

	Init(t)

	tt := time.Now().Unix()
	br := &tables.Order{
		OntId:     "111",
		OrderTime: tt,
	}
	err := TestDB.OrderDB.InsertOrder(br)
	assert.Nil(t, err)

	key := &tables.APIKey{
		ApiKey:  "key",
		Limit:   2,
		UsedNum: 1,
	}
	err = TestDB.ApiDB.InsertApiKey(key)
	assert.Nil(t, err)
	usedNum, err := TestDB.ApiDB.QueryApiKey("key")
	assert.Nil(t, err)
	assert.Equal(t, 1, usedNum)
}

func TestSagaDB_QueryRequestNum(t *testing.T) {
	Init(t)
	usedNum, err := TestDB.ApiDB.QueryApiKey("key")
	assert.Nil(t, err)
	assert.Equal(t, 1, usedNum)
}

func TestSagaDB_SearchApi(t *testing.T) {
	Init(t)
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
	err := TestDB.ApiDB.InsertApiBasicInfo(info)
	assert.Nil(t, err)
	err = TestDB.ApiDB.InsertApiBasicInfo(info2)
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
	Init(t)
	status, err := TestDB.OrderDB.QueryOrderByOrderId("1")
	assert.Nil(t, err)
	fmt.Println("status:", status)
}
