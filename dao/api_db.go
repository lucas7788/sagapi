package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"strings"
)

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

func (this *ApiDB) QueryApiBasicInfoByPage(start, pageSize int) ([]*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info where ApiId limit ?, ?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(start, pageSize)
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

		infos = append(infos, common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc,
			specifications, popularity, delay, successRate, invokeFrequency, createTime))
	}
	return infos, nil
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
func (this *ApiDB) QueryALLApiBasicInfo() ([]*tables.ApiBasicInfo, error) {
	return this.queryApiBasicInfo(false, false, false)
}
func (this *ApiDB) queryApiBasicInfo(newest, hottest, free bool) ([]*tables.ApiBasicInfo, error) {
	var strSql string
	if newest {
		strSql = `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info order by CreateTime limit ?`
	} else if hottest {
		strSql = `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info order by CreateTime limit ?`
	} else if free {
		strSql = `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info where Price='0' limit ?`
	} else {
		strSql = `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info limit ?`
	}

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(sagaconfig.QueryAmt)
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
		res = append(res, common.BuildApiBasicInfo(apiId, apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc,
			specifications, popularity, delay, successRate, invokeFrequency, createTime))
	}
	return res, nil
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

func (this *ApiDB) QueryApiBasicInfoByCategoryId(categoryId int) ([]*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,
Delay,SuccessRate,InvokeFrequency,CreateTime from tbl_api_basic_info where ApiId 
in (select api_id from tbl_api_tag where tag_id=(select id from tbl_tag where category_id=?))`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(categoryId)
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

func (this *ApiDB) QueryApiDetailInfoById(apiId int) (*tables.ApiDetailInfo, error) {
	strSql := "select Id,RequestType,Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario from tbl_api_detail_info where ApiId=?"
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
		var id int
		var requestType, mark, responseParam, responseExample, dataDesc, dataSource, applicationScenario string
		if err = rows.Scan(&id, &requestType, &mark, &responseParam, &responseExample, &dataDesc, &dataSource, &applicationScenario); err != nil {
			return nil, err
		}
		return &tables.ApiDetailInfo{
			Id:                  id,
			RequestType:         requestType,
			ApiId:               apiId,
			Mark:                mark,
			ResponseParam:       responseParam,
			ResponseExample:     responseExample,
			DataDesc:            dataDesc,
			DataSource:          dataSource,
			ApplicationScenario: applicationScenario,
		}, nil
	}
	return nil, fmt.Errorf("not found")
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

func (this *ApiDB) QueryRequestParamByApiDetailInfoId(apiDetailInfoId int) ([]*tables.RequestParam, error) {
	strSql := "select ParamName, ParamType,Required, Note,ValueDesc from tbl_request_param where ApiDetailInfoId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiDetailInfoId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]*tables.RequestParam, 0)
	for rows.Next() {
		var paramName, paramType, note, valueDesc string
		var required bool
		if err = rows.Scan(&paramName, &paramType, &required, &note, &valueDesc); err != nil {
			return nil, err
		}
		rp := &tables.RequestParam{
			ApiDetailInfoId: apiDetailInfoId,
			ParamName:       paramName,
			ParamType:       paramType,
			Required:        required,
			Note:            note,
			ValueDesc:       valueDesc,
		}
		res = append(res, rp)
	}
	return res, nil
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

func (this *ApiDB) QuerySpecificationsBySpecificationsId(id int) (*tables.Specifications, error) {
	strSql := "select Id, Price, Amount from tbl_specifications where Id=?"
	res, err := this.querySpecificationsById(strSql, id)
	if err != nil {
		return nil, err
	}
	if res == nil || len(res) == 0 {
		return nil, fmt.Errorf("no specifications")
	}
	return res[0], nil
}

func (this *ApiDB) QuerySpecificationsByApiDetailId(id int) ([]*tables.Specifications, error) {
	strSql := "select Id, Price, Amount from tbl_specifications where ApiDetailInfoId=?"
	return this.querySpecificationsById(strSql, id)
}

func (this *ApiDB) querySpecificationsById(strSql string, id int) ([]*tables.Specifications, error) {
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(id)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]*tables.Specifications, 0)
	for rows.Next() {
		var amount, specId int
		var price string
		if err = rows.Scan(&specId, &price, &amount); err != nil {
			return nil, err
		}
		rp := &tables.Specifications{
			Id:              specId,
			ApiDetailInfoId: id,
			Price:           price,
			Amount:          amount,
		}
		res = append(res, rp)
	}
	return res, nil
}

func (this *ApiDB) QueryErrorCodeByApiDetailInfoId(apiDetailInfoId int) ([]*tables.ErrorCode, error) {
	strSql := "select ErrorCode, ErrorDesc from tbl_error_code where ApiDetailInfoId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiDetailInfoId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]*tables.ErrorCode, 0)
	for rows.Next() {
		var errorCode int
		var errorDesc string
		if err = rows.Scan(&errorCode, &errorDesc); err != nil {
			return nil, err
		}
		rp := &tables.ErrorCode{
			ApiDetailInfoId: apiDetailInfoId,
			ErrorCode:       errorCode,
			ErrorDesc:       errorDesc,
		}
		res = append(res, rp)
	}
	return res, nil
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

func (this *ApiDB) QueryApiTestKeyByOntIdAndApiId(ontId string, apiId int) (*tables.APIKey, error) {
	strSql := "select ApiKey, RequestLimit, UsedNum from tbl_api_test_key where OntId=? and ApiId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(ontId, apiId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var apiKey string
		var limit, usedNum int
		if err = rows.Scan(&apiKey, &limit, &usedNum); err != nil {
			return nil, err
		}
		return &tables.APIKey{
			ApiKey:       apiKey,
			ApiId:        apiId,
			RequestLimit: limit,
			UsedNum:      usedNum,
			OntId:        ontId,
		}, nil
	}
	return nil, fmt.Errorf("apikey not found")
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

func (this *ApiDB) UpdateApiKeyInvokeFre(apiKey string, usedNum, apiId, invokeFre int) error {
	var strSql string
	if common.IsTestKey(apiKey) {
		strSql = "update tbl_api_test_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	} else {
		strSql = "update tbl_api_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
	}

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(usedNum, invokeFre, apiKey, apiId)
	return err
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

func (this *ApiDB) InsertApiTag(apiTag *tables.ApiTag) error {
	strSql := `insert into tbl_api_tag (id, api_id, tag_id, state, create_time) values (?,?,?,?,?)`

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(apiTag.Id, apiTag.ApiId, apiTag.TagId, apiTag.State, apiTag.CreateTime)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApiDB) InsertTag(tag *tables.Tag) error {
	strSql := `insert into tbl_tag (id, name, category_id, state, create_time) values (?,?,?,?,?)`

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tag.Id, tag.Name, tag.CategoryId, tag.State, tag.CreateTime)
	if err != nil {
		return err
	}

	return nil
}

func (this *ApiDB) InsertCategory(category *tables.Category) error {
	strSql := `insert into tbl_category (id, name_zh, name_en, icon, state) values (?,?,?,?,?)`

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category.Id, category.NameZh, category.NameEn, category.Icon, category.State)
	if err != nil {
		return err
	}

	return nil
}
