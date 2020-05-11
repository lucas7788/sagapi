package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models/tables"
	"strings"
)

func IsErrNoRows(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func (this *SagaApiDB) ClearApiBasicDB() error {
	strSql := "delete from tbl_api_basic_info"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *SagaApiDB) ClearRequestParamDB() error {
	strSql := "delete from tbl_request_param"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *SagaApiDB) ClearApiDetailDB() error {
	strSql := "delete from tbl_api_detail_info"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *SagaApiDB) ClearSpecificationsDB() error {
	strSql := "delete from tbl_specifications"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *SagaApiDB) ClearApiKeyDB() error {
	strSql := "delete from tbl_api_key"
	_, err := this.DB.Exec(strSql)
	if err != nil {
		return err
	}
	strSql2 := "delete from tbl_api_test_key"
	_, err = this.DB.Exec(strSql2)
	return err
}

func (this *SagaApiDB) ClearAll() error {
	strSql := "delete from tbl_api_test_key"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_qr_code"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_key"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_order"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_error_code"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_request_param"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_specifications"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_tag"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_tag"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_category"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_basic_info"
	this.DB.Exec(strSql)
	return nil
}

func (this *SagaApiDB) InsertApiBasicInfo(tx *sqlx.Tx, infos []*tables.ApiBasicInfo) error {
	var err error
	if len(infos) == 0 {
		return nil
	}

	sqlStrArr := make([]string, len(infos))
	for i, info := range infos {
		sqlStrArr[i] = fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%d','%d','%d','%d','%d','%d','%s','%s','%s','%s','%s','%s','%s')",
			info.Coin, info.ApiType, info.Icon, info.Title, info.ApiProvider, info.ApiSagaUrlKey, info.ApiUrl, info.Price,
			info.ApiDesc, info.ErrorDesc, info.Specifications, info.Popularity, info.Delay, info.SuccessRate, info.InvokeFrequency, info.ApiState, info.RequestType, info.Mark, info.ResponseParam, info.ResponseExample, info.DataDesc, info.DataSource, info.ApplicationScenario)
	}
	strSql := `insert into tbl_api_basic_info (Coin,ApiType,Icon,Title,ApiProvider,ApiSagaUrlKey,ApiUrl,Price,ApiDesc,ErrorDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType,Mark,ResponseParam,ResponseExample,DataDesc,DataSource,ApplicationScenario) values`
	strSql += strings.Join(sqlStrArr, ",")
	err = this.Exec(tx, strSql)
	return err
}

func (this *SagaApiDB) QueryApiBasicInfoByApiId(tx *sqlx.Tx, apiId uint32) (*tables.ApiBasicInfo, error) {
	var err error
	strSql := `select * from tbl_api_basic_info where ApiId=? and ApiState=?`

	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, apiId, tables.API_STATE_BUILTIN)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *SagaApiDB) QueryApiBasicInfoBySagaUrlKey(tx *sqlx.Tx, urlkey string) (*tables.ApiBasicInfo, error) {
	var err error
	strSql := `select * from tbl_api_basic_info where ApiSagaUrlKey=? and ApiState=?`
	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, urlkey, tables.API_STATE_BUILTIN)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *SagaApiDB) SearchApi(tx *sqlx.Tx) (map[string][]*tables.ApiBasicInfo, error) {
	res := make(map[string][]*tables.ApiBasicInfo)
	var newestApi []*tables.ApiBasicInfo
	var hottestApi []*tables.ApiBasicInfo
	var freeApi []*tables.ApiBasicInfo
	strNew := "select * from tbl_api_basic_info where ApiState=? order by CreateTime limit ?"
	strHot := "select * from tbl_api_basic_info where ApiState=? order by InvokeFrequency limit ?"
	strFree := "select * from tbl_api_basic_info where Price='0' and ApiState=? limit ?"
	err := this.Select(tx, &newestApi, strNew, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}
	res["newest"] = newestApi

	err = this.Select(tx, &hottestApi, strHot, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}
	res["hottest"] = hottestApi

	err = this.Select(tx, &freeApi, strFree, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}

	res["free"] = freeApi
	return res, nil
}

func (this *SagaApiDB) QueryApiBasicInfoByCategoryId(tx *sqlx.Tx, categoryId, start, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiState=? and ApiId in (select ApiId from tbl_api_tag where TagId=(select Id from tbl_tag where CategoryId=?)) limit ?, ?`

	var res []*tables.ApiBasicInfo
	err := this.Select(tx, &res, strSql, tables.API_STATE_BUILTIN, categoryId, start, pageSize)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *SagaApiDB) QueryApiBasicInfoByPage(start, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiState=? and ApiId limit ?, ?`
	var res []*tables.ApiBasicInfo
	err := this.DB.Select(&res, strSql, tables.API_STATE_BUILTIN, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *SagaApiDB) SearchApiByKey(key string) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select * from tbl_api_basic_info where ApiDesc like ? or Title like ? or ApiId in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where name=?)) limit 30`

	var infos []*tables.ApiBasicInfo
	err := this.DB.Select(&infos, strSql, k, k, key)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (this *SagaApiDB) InsertRequestParam(tx *sqlx.Tx, params []*tables.RequestParam) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		var require int32
		if param.Required {
			require = 1
		} else {
			require = 0
		}
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%d','%d','%s','%s','%s')", param.ApiId, param.ParamName, require, param.ParamWhere, param.ParamType, param.Note, param.ValueDesc)
	}
	strSql := `insert into tbl_request_param (ApiId,ParamName,Required,ParamWhere,ParamType,Note,ValueDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

func (this *SagaApiDB) QueryRequestParamByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.RequestParam, error) {
	strSql := `select * from tbl_request_param where ApiId=?`
	var params []*tables.RequestParam
	err := this.Select(tx, &params, strSql, apiId)
	if err != nil {
		return nil, err
	}
	return params, nil
}

// unit test none.
func (this *SagaApiDB) InsertErrorCode(tx *sqlx.Tx, params []*tables.ErrorCode) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%s')",
			param.ErrorCode, param.ErrorDesc)
	}
	strSql := `insert into tbl_error_code (ErrorCode,ErrorDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

// unit test none.
func (this *SagaApiDB) QueryErrorCode(tx *sqlx.Tx) ([]*tables.ErrorCode, error) {
	strSql := `select * from tbl_error_code`
	var params []*tables.ErrorCode
	err := this.Select(tx, &params, strSql)
	return params, err
}

func (this *SagaApiDB) InsertSpecifications(tx *sqlx.Tx, params []*tables.Specifications) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%d')",
			param.ApiId, param.Price, param.Amount)
	}
	strSql := `insert into tbl_specifications (ApiId,Price,Amount) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

func (this *SagaApiDB) QuerySpecificationsById(tx *sqlx.Tx, id uint32) (*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where Id=?`
	ss := &tables.Specifications{}
	err := this.Get(tx, ss, strSql, id)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (this *SagaApiDB) QuerySpecificationsByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where ApiId=?`
	var ss []*tables.Specifications
	err := this.Select(tx, &ss, strSql, apiId)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

//dependent on orderId.
func (this *SagaApiDB) InsertApiKey(tx *sqlx.Tx, key *tables.APIKey) error {
	strSql := `insert into tbl_api_key (ApiKey,OrderId, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, key.ApiKey, key.OrderId, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	return err
}

//dependent on orderId. use default.
func (this *SagaApiDB) InsertApiTestKey(tx *sqlx.Tx, key *tables.APIKey) error {
	strSql := `insert into tbl_api_test_key (ApiKey, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?)`
	err := this.Exec(tx, strSql, key.ApiKey, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	return err
}

func (this *SagaApiDB) QueryInvokeFreByApiId(tx *sqlx.Tx, apiId uint32) (uint64, error) {
	var freq uint64
	strSql := `select InvokeFrequency from tbl_api_basic_info where ApiId =?`
	err := this.Get(tx, freq, strSql, apiId)
	if err != nil {
		return 0, err
	}

	return freq, nil
}

func (this *SagaApiDB) QueryApiKeyByApiKey(tx *sqlx.Tx, apiKey string) (*tables.APIKey, error) {
	return this.queryApiKey(tx, apiKey, "")
}
func (this *SagaApiDB) QueryApiKeyByOrderId(tx *sqlx.Tx, orderId string) (*tables.APIKey, error) {
	return this.queryApiKey(tx, "", orderId)
}

func (this *SagaApiDB) QueryApiTestKeyByOntidAndApiId(tx *sqlx.Tx, ontid string, apiId uint32) (*tables.APIKey, error) {
	strSql := "select * from tbl_api_test_key where OntId=? and ApiId=?"
	key := &tables.APIKey{}
	err := this.Get(tx, key, strSql, ontid, apiId)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (this *SagaApiDB) queryApiKey(tx *sqlx.Tx, key, orderId string) (*tables.APIKey, error) {
	var strSql string
	var where string
	if key != "" {
		if common.IsTestKey(key) {
			strSql = "select * from tbl_api_test_key where ApiKey=?"
		} else {
			strSql = "select * from tbl_api_key where ApiKey=?"
		}

		where = key
	} else if orderId != "" {
		strSql = "select * from tbl_api_key where OrderId=?"
		where = orderId
	}

	k := &tables.APIKey{}
	err := this.Get(tx, k, strSql, where)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (this *SagaApiDB) UpdateApiKeyInvokeFre(tx *sqlx.Tx, apiKey string, apiId uint32, usedNum, invokeFre uint64) error {
	var strSql string
	var err error
	if common.IsTestKey(apiKey) {
		// here no need update invokefreq.
		strSql = "update tbl_api_test_key set UsedNum=? where ApiKey=?"
		err = this.Exec(tx, strSql, usedNum, apiKey)
	} else {
		strSql = "update tbl_api_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
		err = this.Exec(tx, strSql, usedNum, invokeFre, apiKey, apiId)
	}

	return err
}
