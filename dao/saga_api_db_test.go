package dao

import (
	"database/sql"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiDB_TmpInsert(t *testing.T) {
	TestDB, err := NewSagaApiDB(sagaconfig.DefDBConfigMap[sagaconfig.NETWORK_ID_TRAVIS_NET])
	assert.Nil(t, err)
	tx, err := TestDB.DB.Beginx()
	assert.Nil(t, err)

	TestDB.ClearAll()

	Coin := "ONG"
	ApiType := "ApiType"
	Icon := "Icon"
	Title := "Title"
	ApiProvider := "ApiProvider"
	ApiSagaUrlKey := "ApiSagaUrlKey"
	ApiUrl := "ApiUrl"
	Price := "Price"
	ApiDesc := "ApiDesc"
	ErrorDesc := "ErrorDesc"
	Specifications := uint32(9)
	Popularity := uint32(10)
	Delay := uint32(11)
	SuccessRate := uint32(100)
	InvokeFrequency := uint64(12345)
	ApiState := tables.API_STATE_BUILTIN
	RequestType := "RequestType"
	Mark := "Mark"
	ResponseParam := "ResponseParam"
	ResponseExample := "ResponseExample"
	DataDesc := "DataDesc"
	DataSource := "DataSource"
	ApplicationScenario := "ApplicationScenario"

	info2 := &tables.ApiBasicInfo{
		Coin:                Coin,
		ApiType:             ApiType,
		Icon:                Icon,
		Title:               Title,
		ApiProvider:         ApiProvider,
		ApiSagaUrlKey:       ApiSagaUrlKey,
		ApiUrl:              ApiUrl,
		Price:               Price,
		ApiDesc:             ApiDesc,
		ErrorDesc:           ErrorDesc,
		Specifications:      Specifications,
		Popularity:          Popularity,
		Delay:               Delay,
		SuccessRate:         SuccessRate,
		InvokeFrequency:     InvokeFrequency,
		ApiState:            ApiState,
		RequestType:         RequestType,
		Mark:                Mark,
		ResponseParam:       ResponseParam,
		ResponseExample:     ResponseExample,
		DataDesc:            DataDesc,
		DataSource:          DataSource,
		ApplicationScenario: ApplicationScenario,
	}

	// insert.
	err = TestDB.InsertApiBasicInfo(tx, []*tables.ApiBasicInfo{info2})
	assert.Nil(t, err)

	// try query with tx.
	infoResult, err := TestDB.QueryApiBasicInfoBySagaUrlKey(tx, info2.ApiSagaUrlKey)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	infoResult, err = TestDB.QueryApiBasicInfoByApiId(tx, infoResult.ApiId)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	// try query with db.
	infoResult, err = TestDB.QueryApiBasicInfoBySagaUrlKey(nil, info2.ApiSagaUrlKey)
	assert.Equal(t, err, sql.ErrNoRows)

	err = tx.Commit()
	assert.Nil(t, err)

	// try query with db again.
	infoResult, err = TestDB.QueryApiBasicInfoBySagaUrlKey(nil, info2.ApiSagaUrlKey)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	infoResult, err = TestDB.QueryApiBasicInfoByApiId(nil, infoResult.ApiId)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	l := 11
	infos := make([]*tables.ApiBasicInfo, l)
	for i := 0; i < len(infos); i++ {
		info := &tables.ApiBasicInfo{
			Icon:                "",
			Title:               "mytestasd",
			ApiProvider:         common.GenerateUUId(1),
			ApiSagaUrlKey:       common.GenerateUUId(1),
			ApiUrl:              "",
			Price:               "",
			ApiDesc:             "",
			Specifications:      1,
			ApiState:            tables.API_STATE_BUILTIN,
			Popularity:          0,
			Delay:               0,
			SuccessRate:         0,
			InvokeFrequency:     0,
			ApplicationScenario: common.GenerateUUId(1),
		}

		info.ApiProvider = common.GenerateUUId(1)
		infos[i] = info
	}
	err = TestDB.InsertApiBasicInfo(nil, infos)
	assert.Nil(t, err)

	for i := 0; i < len(infos); i++ {
		infoResult, err := TestDB.QueryApiBasicInfoBySagaUrlKey(nil, infos[i].ApiSagaUrlKey)
		assert.Nil(t, err)
		assert.Equal(t, infoResult.ApplicationScenario, infos[i].ApplicationScenario)
	}

	tx, err = TestDB.DB.Beginx()
	assert.Nil(t, err)
	// test SearchApi
	res, err := TestDB.SearchApi(tx)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(res["newest"]))
	assert.Equal(t, 10, len(res["hottest"]))

	// test RequestParam
	infoResult, err = TestDB.QueryApiBasicInfoBySagaUrlKey(tx, info2.ApiSagaUrlKey)
	assert.Nil(t, err)
	params := &tables.RequestParam{
		ApiId:      infoResult.ApiId,
		ParamName:  "",
		Required:   true,
		ParamType:  "",
		ParamWhere: tables.URL_PARAM_RESTFUL,
		Note:       "",
		ValueDesc:  "ValueDesc",
	}

	assert.Nil(t, TestDB.InsertRequestParam(tx, []*tables.RequestParam{params}))
	paramResult, err := TestDB.QueryRequestParamByApiId(tx, infoResult.ApiId)
	assert.Nil(t, err)
	assert.Equal(t, params.ValueDesc, paramResult[0].ValueDesc)
	assert.Equal(t, params.Required, paramResult[0].Required)

	//TestDB.InsertTag

	err = tx.Commit()
	assert.Nil(t, err)

	err = TestDB.ClearRequestParamDB()
	err = TestDB.ClearApiBasicDB()
	assert.Nil(t, err)
}
