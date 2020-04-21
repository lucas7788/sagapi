package dao

import (
	"github.com/ontio/sagapi/models/tables"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func TestApiDB_InsertApiBasicInfo(t *testing.T) {
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "",
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

	infos, err := TestDB.ApiDB.QueryApiBasicInfoByPage(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, len(infos), 2)
	price, err := TestDB.ApiDB.QueryPriceByApiId(1)
	assert.Nil(t, err)
	assert.Equal(t, price, "1")
}

func TestApiDB_InsertApiDetailInfo(t *testing.T) {
	info := &tables.ApiDetailInfo{
		ApiId:               1,
		Mark:                "",
		ResponseParam:       "",
		ResponseExample:     "",
		DataDesc:            "test",
		DataSource:          "",
		ApplicationScenario: "",
	}
	err := TestDB.ApiDB.InsertApiDetailInfo(info)
	assert.Nil(t, err)
}

func TestApiDB_QueryApiDetailInfoById(t *testing.T) {
	info, err := TestDB.ApiDB.QueryApiDetailInfoById(10)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, 10)
}

func TestApiDB_InsertRequestParam(t *testing.T) {
	rp := &tables.RequestParam{
		ApiDetailInfoId: 1,
		ParamName:       "",
		Required:        true,
		ParamType:       "",
		Note:            "",
	}
	l := 10
	requestParam := make([]*tables.RequestParam, l)
	for i := 0; i < l; i++ {
		requestParam[i] = rp
	}
	err := TestDB.ApiDB.InsertRequestParam(requestParam)
	assert.Nil(t, err)
}

func TestApiDB_QueryRequestParamByApiDetailInfoId(t *testing.T) {
	param, err := TestDB.ApiDB.QueryRequestParamByApiDetailInfoId(1)
	assert.Nil(t, err)
	assert.Equal(t, len(param), 10)
}

func TestApiDB_InsertErrorCode(t *testing.T) {
	rp := &tables.ErrorCode{
		ApiDetailInfoId: 1,
		ErrorCode:       1,
		ErrorDesc:       "",
	}
	l := 10
	requestParam := make([]*tables.ErrorCode, l)
	for i := 0; i < l; i++ {
		requestParam[i] = rp
	}
	err := TestDB.ApiDB.InsertErrorCode(requestParam)
	assert.Nil(t, err)
}

func TestApiDB_QueryErrorCodeByApiDetailInfoId(t *testing.T) {
	param, err := TestDB.ApiDB.QueryErrorCodeByApiDetailInfoId(1)
	assert.Nil(t, err)
	assert.Equal(t, len(param), 10)
}

func TestApiDB_QueryNewestApiBasicInfo(t *testing.T) {
	infos, err := TestDB.ApiDB.QueryFreeApiBasicInfo()
	fmt.Println(err)
	assert.Nil(t,err)
	fmt.Println(infos)
}