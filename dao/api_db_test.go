package dao

import (
	"database/sql"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestApiDB *ApiDB
var TestOrderDB *OrderDB

func InitApiDb(t *testing.T) {
	db, dberr := sql.Open("mysql",
		config.DefConfig.DbConfig.ProjectDBUser+
			":"+config.DefConfig.DbConfig.ProjectDBPassword+
			"@tcp("+config.DefConfig.DbConfig.ProjectDBUrl+
			")/"+config.DefConfig.DbConfig.ProjectDBName+
			"?charset=utf8")
	assert.Nil(t, dberr)
	err := db.Ping()
	assert.Nil(t, err)
	TestApiDB = NewApiDB(db)
	TestOrderDB = NewOrderDB(db)
}

func TestApiDB_InsertApiBasicInfo(t *testing.T) {
	InitApiDb(t)
	info := &tables.ApiBasicInfo{
		ApiLogo:         "",
		ApiName:         "",
		ApiProvider:     "",
		ApiUrl:          "",
		ApiPrice:        "",
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
	assert.Nil(t, TestApiDB.InsertApiBasicInfo(infos))
}

func TestApiDB_QueryApiBasicInfoByApiId(t *testing.T) {
	InitApiDb(t)
	info, err := TestApiDB.QueryApiBasicInfoByApiId(1)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, 1)

	infos, err := TestApiDB.QueryApiBasicInfoByPage(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, len(infos), 2)
	price, err := TestApiDB.QueryPriceByApiId(1)
	assert.Nil(t, err)
	assert.Equal(t, price, "1")
}

func TestApiDB_InsertApiDetailInfo(t *testing.T) {
	InitApiDb(t)
	info := &tables.ApiDetailInfo{
		ApiId:               10,
		Mark:                "",
		ResponseParam:       "",
		ResponseExample:     "",
		DataDesc:            "",
		DataSource:          "",
		ApplicationScenario: "",
	}
	err := TestApiDB.InsertApiDetailInfo(info)
	assert.Nil(t, err)
}

func TestApiDB_QueryApiDetailInfoById(t *testing.T) {
	InitApiDb(t)
	info, err := TestApiDB.QueryApiDetailInfoById(10)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, 10)
}
