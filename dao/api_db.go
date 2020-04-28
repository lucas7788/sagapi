package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"strings"
)

type ApiDB struct {
	conn *sqlx.DB
}

func IsNoEltError(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func NewApiDB(db *sqlx.DB) *ApiDB {
	return &ApiDB{
		conn: db,
	}
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
	_, err := this.conn.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryApiBasicInfoByCategoryId(categoryId, start, pageSize int) ([]*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiId 
in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where category_id=?)) limit ?, ?`

	var res []*tables.ApiBasicInfo
	err := this.conn.Select(&res, strSql, categoryId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *ApiDB) QueryApiBasicInfoByApiId(apiId int) (*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiId =?`

	info := &tables.ApiBasicInfo{}
	err := this.conn.Get(info, strSql, apiId)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (this *ApiDB) SearchApiByKey(key string) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select * from tbl_api_basic_info where ApiDesc like ? or Title like ? or ApiId in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where name=?)) limit 30`

	var infos []*tables.ApiBasicInfo
	err := this.conn.Select(&infos, strSql, k, k, key)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (this *ApiDB) InsertApiDetailInfo(info *tables.ApiDetailInfo) error {
	strSql := `insert into tbl_api_detail_info (ApiId,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, 
DataSource,ApplicationScenario) values (?,?,?,?,?,?,?,?)`

	_, err := this.conn.Exec(strSql, info.ApiId, info.RequestType, info.Mark, info.ResponseParam, info.ResponseExample, info.DataDesc, info.DataSource, info.ApplicationScenario)
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
	_, err := this.conn.Exec(strSql)
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
	_, err := this.conn.Exec(strSql)
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
	_, err := this.conn.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

//dependent on orderId
func (this *ApiDB) InsertApiKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_key (ApiKey,OrderId, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?,?)`

	_, err := this.conn.Exec(strSql, key.ApiKey, key.OrderId, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

//dependent on orderId
func (this *ApiDB) InsertApiTestKey(key *tables.APIKey) error {
	strSql := `insert into tbl_api_test_key (ApiKey, ApiId, RequestLimit, UsedNum, OntId) values (?,?,?,?,?)`
	_, err := this.conn.Exec(strSql, key.ApiKey, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId)
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

	key := &models.ApiKeyInvokeFre{}
	err := this.conn.Get(key, strSql, apiKey)
	if err != nil {
		return nil, err
	}
	return key, nil
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
	k := &tables.APIKey{}
	err := this.conn.Get(k, strSql, where)
	if err != nil {
		return nil, err
	}
	return k, nil
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

func (this *ApiDB) UpdateApiKeyInvokeFre(apiKey string, usedNum, apiId, invokeFre int) error {
	var strSql string
	if common.IsTestKey(apiKey) {
		strSql = "update tbl_api_test_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	} else {
		strSql = "update tbl_api_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	}

	_, err := this.conn.Exec(strSql, usedNum, invokeFre, apiKey, apiId)
	return err
}
