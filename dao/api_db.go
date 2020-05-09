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

func IsNoEltError(err error) bool {
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
	strSql := `select * from tbl_api_basic_info where ApiId =?`

	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, apiId)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *SagaApiDB) QueryApiBasicInfoBySagaUrlKey(tx *sqlx.Tx, urlkey string) (*tables.ApiBasicInfo, error) {
	var err error
	strSql := `select * from tbl_api_basic_info where ApiSagaUrlKey=?`
	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, urlkey)
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
	strSql := `select * from tbl_api_basic_info where ApiId in (select ApiId from tbl_api_tag where TagId=(select Id from tbl_tag where CategoryId=?)) limit ?, ?`

	var res []*tables.ApiBasicInfo
	err := this.Select(tx, &res, strSql, categoryId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *SagaApiDB) QueryApiBasicInfoByPage(start, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiId limit ?, ?`
	var res []*tables.ApiBasicInfo
	err := this.DB.Select(&res, strSql, start, pageSize)
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

func (this *SagaApiDB) InsertErrorCode(params []*tables.ErrorCode) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%d','%s')",
			param.ApiId, param.ErrorCode, param.ErrorDesc)
	}
	strSql := `insert into tbl_error_code (ApiId,ErrorCode,ErrorDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.DB.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *SagaApiDB) QueryErrorCodeByApiDetailId(id uint32) ([]*tables.ErrorCode, error) {
	strSql := `select * from tbl_error_code where ApiId=?`
	var params []*tables.ErrorCode
	err := this.DB.Select(&params, strSql, id)
	return params, err
}

func (this *SagaApiDB) InsertSpecifications(params []*tables.Specifications) error {
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
	_, err := this.DB.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *SagaApiDB) QuerySpecificationsById(id uint32) (*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where Id=?`
	ss := &tables.Specifications{}
	err := this.DB.Get(ss, strSql, id)
	return ss, err
}

func (this *SagaApiDB) QuerySpecificationsByApiDetailId(id uint32) ([]*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where ApiId=?`
	var ss []*tables.Specifications
	err := this.DB.Select(&ss, strSql, id)
	return ss, err
}

//dependent on orderId
func (this *SagaApiDB) InsertApiKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_key (ApiKey,OrderId, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?,?)`

	_, err := this.DB.Exec(strSql, key.ApiKey, key.OrderId, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

//dependent on orderId
func (this *SagaApiDB) InsertApiTestKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_test_key (ApiKey, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?)`
	_, err := this.DB.Exec(strSql, key.ApiKey, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

func (this *SagaApiDB) QueryInvokeFreByApiId(apiId uint32) (uint64, error) {
	var freq uint64
	strSql := `select InvokeFrequency from tbl_api_basic_info where ApiId =?`
	err := this.DB.Get(freq, strSql, apiId)
	if err != nil {
		return 0, err
	}

	return freq, nil
}

func (this *SagaApiDB) QueryApiKeyByApiKey(apiKey string) (*tables.APIKey, error) {
	return this.queryApiKey(apiKey, "")
}
func (this *SagaApiDB) QueryApiKeyByOrderId(orderId string) (*tables.APIKey, error) {
	return this.queryApiKey("", orderId)
}

func (this *SagaApiDB) QueryApiTestKeyByOntidAndApiId(ontid string, apiId uint32) (*tables.APIKey, error) {
	strSql := "select * from tbl_api_test_key where OntId=? and ApiId=?"
	key := &tables.APIKey{}
	err := this.DB.Get(key, strSql, ontid, apiId)
	return key, err
}

func (this *SagaApiDB) queryApiKey(key, orderId string) (*tables.APIKey, error) {
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
	err := this.DB.Get(k, strSql, where)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (this *SagaApiDB) VerifyApiKey(apiKey string) error {
	key, err := this.QueryApiKeyByApiKey(apiKey)
	if err != nil {
		return err
	}
	if key == nil {
		return fmt.Errorf("invalid api key: %s", apiKey)
	}
	if key.UsedNum >= key.RequestLimit {
		return fmt.Errorf("available times:%d, has used times: %d", key.RequestLimit, key.UsedNum)
	}
	return nil
}

func (this *SagaApiDB) UpdateApiKeyInvokeFre(apiKey string, apiId uint32, usedNum, invokeFre uint64) error {
	var strSql string
	if common.IsTestKey(apiKey) {
		strSql = "update tbl_api_test_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	} else {
		strSql = "update tbl_api_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	}

	_, err := this.DB.Exec(strSql, usedNum, invokeFre, apiKey, apiId)
	return err
}
