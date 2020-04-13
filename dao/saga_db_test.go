package dao

import (
	"testing"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/saga/config"
	"github.com/ontio/saga/models/tables"
	"github.com/stretchr/testify/assert"
)

var TestDB *SagaDB

func Init(t *testing.T) {
	config.DefConfig.ProjectDBHost = "127.0.0.1:30336"
	config.DefConfig.ProjectDBName = "saga"
	config.DefConfig.ProjectDBUser = "root"
	config.DefConfig.ProjectDBPassword = "111111"

	db, err := NewDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Init()
	assert.Nil(t, err)
	TestDB = db
}

func TestSagaDB_Init(t *testing.T) {

	Init(t)

	br := &tables.BuyRecord{
		OntId: "111",
	}
	err := TestDB.InsertBuyRecord(br)
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
	//info := &tables.APIInfo{
	//	ApiDesc:"abcdefg",
	//}
	//info2 := &tables.APIInfo{
	//	ApiDesc:"cdefgty",
	//}
	//err := TestDB.InsertApiInfo(info)
	//assert.Nil(t, err)
	//err = TestDB.InsertApiInfo(info2)
	//assert.Nil(t, err)
	infos, err := TestDB.SearchApi("cdefgty")
	assert.Nil(t, err)
	fmt.Println(infos)
	infos, err = TestDB.QueryApiInfoByPage(2, 2)
	assert.Nil(t, err)
	fmt.Println(infos)
	info3, err := TestDB.QueryApiInfoByApiId(100)
	assert.Nil(t, err)
	fmt.Println(info3)
	TestDB.Close()
}
