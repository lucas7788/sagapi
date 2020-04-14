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

var TestDB *SagaDB

func Init(t *testing.T) {
	config.DefConfig.DbConfig.ProjectDBUrl = "127.0.0.1:3306"
	config.DefConfig.DbConfig.ProjectDBName = "saga"
	config.DefConfig.DbConfig.ProjectDBUser = "root"
	config.DefConfig.DbConfig.ProjectDBPassword = "111111"

	db, err := NewDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Init()
	assert.Nil(t, err)
	TestDB = db
}

func TestSagaDB_Init(t *testing.T) {

	Init(t)

	tt := time.Now().Unix()
	br := &tables.Order{
		OntId:     "111",
		OrderTime: tt,
	}
	err := TestDB.InsertOrder(br)
	assert.Nil(t, err)

	key := &tables.APIKey{
		ApiKey:  "key",
		Limit:   2,
		UsedNum: 1,
	}
	err = TestDB.InsertApiKey(key)
	assert.Nil(t, err)
	usedNum, err := TestDB.QueryRequestNum("key")
	assert.Nil(t, err)
	assert.Equal(t, 1, usedNum)
}

func TestSagaDB_QueryRequestNum(t *testing.T) {
	Init(t)
	usedNum, err := TestDB.QueryRequestNum("key")
	assert.Nil(t, err)
	assert.Equal(t, 1, usedNum)
}

func TestSagaDB_SearchApi(t *testing.T) {
	Init(t)
	info := &tables.ApiBasicInfo{
		ApiDesc: "abcdefg",
	}
	info2 := &tables.ApiBasicInfo{
		ApiDesc: "cdefgty",
	}
	err := TestDB.InsertApiBasicInfo(info)
	assert.Nil(t, err)
	err = TestDB.InsertApiBasicInfo(info2)
	assert.Nil(t, err)
	infos, err := TestDB.SearchApi("cdefgty")
	assert.Nil(t, err)
	fmt.Println(infos)
	infos, err = TestDB.QueryApiInfoByPage(2, 2)
	assert.Nil(t, err)
	fmt.Println(infos)
	info3, err := TestDB.QueryApiBasicInfoByApiId(100)
	assert.Nil(t, err)
	fmt.Println(info3)
	TestDB.Close()
}

func TestSagaDB_QueryOrderStatusByOrderId(t *testing.T) {
	Init(t)
	status, err := TestDB.QueryOrderStatusByOrderId("1")
	assert.Nil(t, err)
	fmt.Println("status:", status)
}