package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"strings"
)

var DefApiDb = ApiDb{}

type ApiDb struct {
	Conn *sqlx.DB
}

func InitDefApiDb(db *sqlx.DB) {
	log.Info("InitDefApiDb Done.")
	DefApiDb.Conn = db
}

func IsNoEltError(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

type ApiDB struct {
	db *sql.DB
}

func NewApiDB(db *sql.DB) *ApiDB {
	return &ApiDB{
		db: db,
	}
}

func (this *ApiDB) DB() *sql.DB {
	return this.db
}

func (this *ApiDB) InsertApiBasicInfo(infos []*tables.ApiBasicInfo) error {
	if len(infos) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(infos))
	for i, info := range infos {
		sqlStrArr[i] = fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%d','%d','%d','%d','%d')",
			info.Coin, info.ApiType, info.Icon, info.Title, info.ApiProvider, info.ApiUrl, info.Price,
			info.ApiDesc, info.Specifications, info.Popularity, info.Delay, info.SuccessRate, info.InvokeFrequency)
	}
	strSql := `insert into tbl_api_basic_info (Coin,ApiType,Icon, Title, ApiProvider, ApiUrl, Price, 
ApiDesc,Specifications, Popularity,Delay,SuccessRate,InvokeFrequency) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.db.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryHottestApiBasicInfo() ([]*tables.ApiBasicInfo, error) {
	return this.queryApiBasicInfo(false, true, false)
}
func (this *ApiDB) QueryNewestApiBasicInfo() ([]*tables.ApiBasicInfo, error) {
	return this.queryApiBasicInfo(true, false, false)
}
func (this *ApiDB) QueryFreeApiBasicInfo() ([]*tables.ApiBasicInfo, error) {
	return this.queryApiBasicInfo(false, false, true)
}

func (this *ApiDB) QueryInvokeFreByApiId(apiId int) (int, error) {
	strSql := `select InvokeFrequency from tbl_api_basic_info where ApiId =?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query(apiId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var invokeFrequency int
		if err = rows.Scan(&invokeFrequency); err != nil {
			return 0, err
		}
		return invokeFrequency, nil
	}
	return 0, fmt.Errorf("not found")
}

func (this *ApiDB) QueryApiBasicInfoByCategoryId(categoryId, start, pageSize int) ([]*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,
Delay,SuccessRate,InvokeFrequency,CreateTime from tbl_api_basic_info where ApiId 
in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where category_id=?)) limit ?, ?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(categoryId, start, pageSize)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]*tables.ApiBasicInfo, 0)
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, createTime string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications,
			&popularity, &delay, &successRate, &invokeFrequency, &createTime); err != nil {
			return nil, err
		}
		res = append(res, common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, specifications,
			popularity, delay, successRate, invokeFrequency, createTime))
	}
	return res, nil
}

func (this *ApiDB) QueryApiBasicInfoByApiId(apiId int) (*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,
Delay,SuccessRate,InvokeFrequency,CreateTime from tbl_api_basic_info where ApiId =?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, createTime string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications,
			&popularity, &delay, &successRate, &invokeFrequency, &createTime); err != nil {
			return nil, err
		}
		return common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, specifications,
			popularity, delay, successRate, invokeFrequency, createTime), nil
	}
	return nil, fmt.Errorf("not found")
}

func (this *ApiDB) QueryApiByApiIds(apiIds []int) ([]*tables.ApiBasicInfo, error) {
	res := make([]*tables.ApiBasicInfo, len(apiIds))
	for i, apiId := range apiIds {
		info, err := this.QueryApiBasicInfoByApiId(apiId)
		if err != nil {
			return nil, err
		}
		res[i] = info
	}
	return res, nil
}

func (this *ApiDB) QueryApiByApiId(apiId int) (*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,
SuccessRate,InvokeFrequency,CreateTime from tbl_api_basic_info where ApiId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, createTime string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications,
			&popularity, &delay, &successRate, &invokeFrequency, &createTime); err != nil {
			return nil, err
		}
		return common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, specifications, popularity,
			delay, successRate, invokeFrequency, createTime), nil
	}
	return nil, fmt.Errorf("not found")
}

func (this *ApiDB) SearchApiByKey(key string) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info where ApiDesc like ? or Title like ? or ApiId in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where name=?)) limit 30`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(k, k, key)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	infos := make([]*tables.ApiBasicInfo, 0)
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, createTime string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications,
			&popularity, &delay, &successRate, &invokeFrequency, &createTime); err != nil {
			return nil, err
		}
		infos = append(infos, common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc, specifications, popularity,
			delay, successRate, invokeFrequency, createTime))
	}
	return infos, nil
}

func (this *ApiDB) InsertApiDetailInfo(info *tables.ApiDetailInfo) error {
	strSql := `insert into tbl_api_detail_info (ApiId,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, 
DataSource,ApplicationScenario) values (?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ApiId, info.RequestType, info.Mark, info.ResponseParam, info.ResponseExample, info.DataDesc, info.DataSource, info.ApplicationScenario)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) InsertRequestParam(params []*tables.RequestParam) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		var require int
		if param.Required {
			require = 1
		} else {
			require = 0
		}
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%s','%d','%s')",
			param.ApiDetailInfoId, param.ParamName, param.ParamType, require, param.ValueDesc)
	}
	strSql := `insert into tbl_request_param (ApiDetailInfoId,ParamName,ParamType,Required,ValueDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.db.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) InsertErrorCode(params []*tables.ErrorCode) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%d','%s')",
			param.ApiDetailInfoId, param.ErrorCode, param.ErrorDesc)
	}
	strSql := `insert into tbl_error_code (ApiDetailInfoId,ErrorCode,ErrorDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.db.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}
func (this *ApiDB) InsertSpecifications(params []*tables.Specifications) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%d')",
			param.ApiDetailInfoId, param.Price, param.Amount)
	}
	strSql := `insert into tbl_specifications (ApiDetailInfoId,Price,Amount) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.db.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryPriceByApiId(ApiId int) (string, error) {
	strSql := "select Price from tbl_api_basic_info where ApiId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return "", err
	}
	rows, err := stmt.Query(ApiId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var price string
		if err = rows.Scan(&price); err != nil {
			return "", err
		}
		return price, nil
	}
	return "", nil
}

//dependent on orderId
func (this *ApiDB) InsertApiKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_key (ApiKey,OrderId, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key.ApiKey, key.OrderId, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

//dependent on orderId
func (this *ApiDB) InsertApiTestKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_test_key (ApiKey, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key.ApiKey, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryApiKeyAndInvokeFreByApiKey(apiKey string) (*models.ApiKeyInvokeFre, error) {
	var strSql string
	if common.IsTestKey(apiKey) {
		strSql = `select k.ApiId, k.OntId, k.RequestLimit, k.UsedNum,i.InvokeFrequency from tbl_api_test_key k,
tbl_api_basic_info i where k.ApiKey=? and i.ApiId=k.ApiId`
	} else {
		strSql = `select k.ApiId, k.OntId, k.RequestLimit, k.UsedNum,i.InvokeFrequency from tbl_api_key k,
tbl_api_basic_info i where k.ApiKey=? and i.ApiId=k.ApiId`
	}

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiKey)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ontId string
		var apiId, limit, usedNum, invokeFre int
		if err = rows.Scan(&apiId, &ontId, &limit, &usedNum, &invokeFre); err != nil {
			return nil, err
		}
		return &models.ApiKeyInvokeFre{
			ApiKey:       apiKey,
			ApiId:        apiId,
			RequestLimit: limit,
			UsedNum:      int32(usedNum),
			OntId:        ontId,
		}, nil
	}
	return nil, fmt.Errorf("not found")
}

func (this *ApiDB) QueryApiKeyByApiKey(apiKey string) (*tables.APIKey, error) {
	return this.queryApiKey(apiKey, "")
}
func (this *ApiDB) QueryApiKeyByOrderId(orderId string) (*tables.APIKey, error) {
	return this.queryApiKey("", orderId)
}

func (this *ApiDB) queryApiKey(key, orderId string) (*tables.APIKey, error) {
	var strSql string
	var where string
	if key != "" {
		strSql = "select ApiKey, OrderId, ApiId, RequestLimit, UsedNum, OntId from tbl_api_key where ApiKey=?"
		where = key
	} else if orderId != "" {
		strSql = "select ApiKey, OrderId, ApiId, RequestLimit, UsedNum, OntId from tbl_api_key where OrderId=?"
		where = orderId
	}
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(where)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ontId, orderId, apiKey string
		var apiId, limit, usedNum int
		if err = rows.Scan(&apiKey, &orderId, &apiId, &limit, &usedNum, &ontId); err != nil {
			return nil, err
		}
		return &tables.APIKey{
			OrderId:      orderId,
			ApiKey:       apiKey,
			ApiId:        apiId,
			RequestLimit: limit,
			UsedNum:      usedNum,
			OntId:        ontId,
		}, nil
	}
	return nil, fmt.Errorf("not found")
}

func (this *ApiDB) VerifyApiKey(apiKey string) error {
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
