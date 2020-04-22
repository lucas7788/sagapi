package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
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

func (this *ApiDB) UpdateInvokeFrequencyByApiId(invokeFre, apiId int) error {
	strSql := `update tbl_api_basic_info set InvokeFrequency=? where ApiId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(invokeFre, apiId)
	return err
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
InvokeFrequency,CreateTime from tbl_api_basic_info where Price='' limit ?`
	}

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(config.QueryAmt)
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
	return 0, nil
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
	return nil, nil
}
func (this *ApiDB) QueryTagIdByCategoryId(categoryId int) (int, error) {
	strSql := `select id from tbl_tag where category_id =?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query(categoryId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var tagId int
		if err = rows.Scan(&tagId); err != nil {
			return 0, err
		}
		return tagId, nil
	}
	return 0, nil
}
func (this *ApiDB) QueryApiIdByTagId(tagId int) ([]int, error) {
	strSql := `select api_id from tbl_api_tag where tag_id =?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(tagId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]int, 0)
	for rows.Next() {
		var apiId int
		if err = rows.Scan(&apiId); err != nil {
			return nil, err
		}
		res = append(res, apiId)
	}
	return res, nil
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
	return nil, nil
}

func (this *ApiDB) SearchApiByKey(key string) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select ApiId, Icon, Title, ApiProvider, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency,CreateTime from tbl_api_basic_info where ApiDesc like ? or Title like ? limit 10`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(k, k)
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
	return nil, nil
}

func (this *ApiDB) InsertRequestParam(params []*tables.RequestParam) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%s','%t','%s')",
			param.ApiDetailInfoId, param.ParamName, param.ParamType, param.Required, param.Note)
	}
	strSql := `insert into tbl_request_param (ApiDetailInfoId,ParamName,ParamType,Required,Note) values`
	strSql += strings.Join(sqlStrArr, ",")
	_, err := this.db.Exec(strSql)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryRequestParamByApiDetailInfoId(apiDetailInfoId int) ([]*tables.RequestParam, error) {
	strSql := "select ParamName, ParamType,Required, Note from tbl_request_param where ApiDetailInfoId=?"
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
		var paramName, paramType, note string
		var required bool
		if err = rows.Scan(&paramName, &paramType, &required, &note); err != nil {
			return nil, err
		}
		rp := &tables.RequestParam{
			ApiDetailInfoId: apiDetailInfoId,
			ParamName:       paramName,
			ParamType:       paramType,
			Required:        required,
			Note:            note,
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
	return nil, nil
}

func (this *ApiDB) QueryApiTestKeyByApiTestKey(apiTestKey string) (*tables.APIKey, error) {
	strSql := "select ApiId, OntId, RequestLimit, UsedNum from tbl_api_test_key where ApiKey=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(apiTestKey)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ontId string
		var apiId, limit, usedNum int
		if err = rows.Scan(&apiId, &ontId, &limit, &usedNum); err != nil {
			return nil, err
		}
		return &tables.APIKey{
			ApiKey:       apiTestKey,
			ApiId:        apiId,
			RequestLimit: limit,
			UsedNum:      usedNum,
			OntId:        ontId,
		}, nil
	}
	return nil, nil
}

func (this *ApiDB) UpdateApiKeyUsedNum(apiKey string, usedNum int) error {
	strSql := "update tbl_api_key set UsedNum = ? where ApiKey=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(usedNum, apiKey)
	return err
}

func (this *ApiDB) UpdateApiTestKeyUsedNum(apiKey string, usedNum int) error {
	strSql := "update tbl_api_test_key set UsedNum = ? where ApiKey=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(usedNum, apiKey)
	return err
}

func (this *ApiDB) QueryApiKey(apiKey string) (*tables.APIKey, error) {
	strSql := "select OrderId, ApiId, RequestLimit, UsedNum, OntId from tbl_api_key where ApiKey=?"
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
		var ontId, orderId string
		var apiId, limit, usedNum int
		if err = rows.Scan(&orderId, &apiId, &limit, &usedNum, &ontId); err != nil {
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
	return nil, nil
}

func (this *ApiDB) VerifyApiKey(apiKey string) error {
	key, err := this.QueryApiKey(apiKey)
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
	strSql := `insert into tbl_tag (id, name_zh, name_en, icon, state) values (?,?,?,?,?)`

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
