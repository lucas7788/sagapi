package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/models/tables"
)

type ApiDB struct {
	db *sql.DB
}

func NewApiDB(db *sql.DB) *ApiDB {
	return &ApiDB{
		db: db,
	}
}

func (this *ApiDB) InsertApiBasicInfo(info *tables.ApiBasicInfo) error {
	strSql := `insert into api_basic_info (ApiLogo, ApiName, ApiProvider, ApiUrl, ApiPrice, ApiDesc,Specifications, 
Popularity,Delay,SuccessRate,InvokeFrequency) values (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ApiLogo, info.ApiName, info.ApiProvider, info.ApiUrl, info.ApiPrice, info.ApiDesc,
		info.Specifications, info.Popularity, info.Delay, info.SuccessRate, info.InvokeFrequency)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryApiBasicInfoByPage(start, pageSize int) (infos []*tables.ApiBasicInfo, err error) {
	strSql := `select ApiId, ApiLogo, ApiName, ApiProvider, ApiUrl, ApiPrice, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency from api_basic_info_tbl where order by id limit ?, ?`
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
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications, &popularity, &delay, &successRate, &invokeFrequency); err != nil {
			return nil, err
		}
		infos = append(infos, &tables.ApiBasicInfo{
			ApiId:           apiId,
			ApiLogo:         apiLogo,
			ApiName:         apiName,
			ApiProvider:     apiProvider,
			ApiUrl:          apiUrl,
			ApiPrice:        apiPrice,
			ApiDesc:         apiDesc,
			Specifications:  specifications,
			Popularity:      popularity,
			Delay:           delay,
			SuccessRate:     successRate,
			InvokeFrequency: invokeFrequency,
		})
	}
	return
}

func (this *ApiDB) QueryApiBasicInfoByApiId(apiId int) (*tables.ApiBasicInfo, error) {
	strSql := `select ApiId, ApiLogo, ApiName, ApiProvider, ApiUrl, ApiPrice, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency from api_basic_info_tbl where order by ApiId =?`
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
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications, &popularity, &delay, &successRate, &invokeFrequency); err != nil {
			return nil, err
		}
		return &tables.ApiBasicInfo{
			ApiId:           apiId,
			ApiLogo:         apiLogo,
			ApiName:         apiName,
			ApiProvider:     apiProvider,
			ApiUrl:          apiUrl,
			ApiPrice:        apiPrice,
			ApiDesc:         apiDesc,
			Specifications:  specifications,
			Popularity:      popularity,
			Delay:           delay,
			SuccessRate:     successRate,
			InvokeFrequency: invokeFrequency,
		}, nil
	}
	return nil, nil
}

func (this *ApiDB) SearchApi(key string) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select ApiId, ApiLogo, ApiName, ApiProvider, ApiUrl, ApiPrice, ApiDesc,Specifications,Popularity,Delay,SuccessRate,
InvokeFrequency from api_basic_info where api_desc like ?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(k)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	infos := make([]*tables.ApiBasicInfo, 0)
	for rows.Next() {
		var apiLogo, apiName, apiProvider, apiUrl, apiPrice, apiDesc string
		var apiId, specifications, popularity, delay, successRate, invokeFrequency int
		if err = rows.Scan(&apiId, &apiLogo, &apiName, &apiProvider, &apiUrl, &apiPrice, &apiDesc, &specifications, &popularity, &delay, &successRate, &invokeFrequency); err != nil {
			return nil, err
		}
		infos = append(infos, &tables.ApiBasicInfo{
			ApiId:           apiId,
			ApiLogo:         apiLogo,
			ApiName:         apiName,
			ApiProvider:     apiProvider,
			ApiUrl:          apiUrl,
			ApiPrice:        apiPrice,
			ApiDesc:         apiDesc,
			Specifications:  specifications,
			Popularity:      popularity,
			Delay:           delay,
			SuccessRate:     successRate,
			InvokeFrequency: invokeFrequency,
		})
	}
	return infos, nil
}

func (this *ApiDB) InsertApiDetailInfo(info *tables.ApiDetailInfo) error {
	strSql := `insert into api_detail_info (ApiId, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario) values (?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ApiId, info.Mark, info.ResponseParam, info.ResponseExample, info.DataDesc, info.DataSource, info.ApplicationScenario)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) QueryApiDetailInfoById(apiId int) (*tables.ApiDetailInfo, error) {
	strSql := "select Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario from api_detail_info where ApiId=?"
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
		var mark, responseParam, responseExample, dataDesc, dataSource, applicationScenario string
		if err = rows.Scan(&mark, &responseParam, &responseExample, &dataDesc, &dataSource, &applicationScenario); err != nil {
			return nil, err
		}
		return &tables.ApiDetailInfo{
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

func (this *ApiDB) QueryPriceByApiId(ApiId int) (string, error) {
	strSql := "select ApiPrice from api_basic_info where id=?"
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

func (this *ApiDB) InsertApiKey(key *tables.APIKey) error {
	strSql := `insert into api_key (ApiKey, ApiId, Limit, UsedNum, OntId) values (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(key.ApiKey, key.ApiId, key.Limit, key.UsedNum, key.OntId)
	if err != nil {
		return err
	}
	return nil
}

func (this *ApiDB) UpdateApiKeyUsedNum(key string, usedNum int) error {
	strSql := "update api_key_tbl set UsedNum = ? where api_key=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(usedNum, key)
	return err
}

func (this *ApiDB) QueryApiKey(apiKey string) (*tables.APIKey, error) {
	strSql := "select ApiId, Limit, UsedNum, OntId from api_key_tbl where ApiKey=?"
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
		var apiId, limit, usedNum int
		if err = rows.Scan(&apiId, &limit, &usedNum, &ontId); err != nil {
			return nil, err
		}
		return &tables.APIKey{
			ApiKey:  apiKey,
			ApiId:   apiId,
			Limit:   limit,
			UsedNum: usedNum,
			OntId:   ontId,
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
	if key.UsedNum >= key.Limit {
		return fmt.Errorf("Available times:%d, has used times: %d", key.Limit, key.UsedNum)
	}
	return nil
}
