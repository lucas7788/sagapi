package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/models/tables"
)

func (this *SagaApiDB) QueryLocationOfCountryCity(tx *sqlx.Tx, country string) ([]*tables.Location, error) {
	var err error
	res := make([]*tables.Location, 0)
	if country != "ALL" {
		strSql := `select * from tbl_country_city where Country=?`

		err = this.Select(tx, &res, strSql, country)
		if err != nil {
			return nil, err
		}
	} else {
		strSql := `select DISTINCT Country from tbl_country_city`

		err = this.Select(tx, &res, strSql)
		if err != nil {
			return nil, err
		}

	}

	return res, nil
}

func (this *SagaApiDB) QueryApiBasicInfoByApiTypeKind(tx *sqlx.Tx, apiType string, apiKind int32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	var err error
	StrSql := `select * from tbl_api_basic_info where ApiType=? and ApiKind=? and ApiState=?`
	res := make([]*tables.ApiBasicInfo, 0)
	err = this.Select(tx, &res, StrSql, apiType, apiKind, apiState)
	if err != nil {
		return nil, err
	}
	log.Debugf("apiType: %s, %d, %d, len(res): %d", apiType, apiKind, apiState, len(res))

	return res, nil
}

func (this *SagaApiDB) QueryApiAlgorithmsByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.ApiAlgorithm, error) {
	var err error
	StrSql := `select * from tbl_api_algorithm where ApiId=? and State=1`
	res := make([]*tables.ApiAlgorithm, 0)
	err = this.Select(tx, &res, StrSql, apiId)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (this *SagaApiDB) QueryAlgorithmEnvByAlgorithmId(tx *sqlx.Tx, algorithmId uint32) ([]*tables.AlgorithmEnv, error) {
	var err error
	StrSql := `select * from tbl_algorithm_env where AlgorithmId=? and State=1`
	res := make([]*tables.AlgorithmEnv, 0)
	err = this.Select(tx, &res, StrSql, algorithmId)
	if err != nil {
		return nil, err
	}
	log.Debugf("algorithmId: %d, len(res):%d", algorithmId, len(res))

	return res, nil

}

func (this *SagaApiDB) QueryAlgorithmsById(tx *sqlx.Tx, id uint32) ([]*tables.Algorithm, error) {
	var err error
	res := make([]*tables.Algorithm, 0)
	if id != 0 {
		StrSql := `select * from tbl_algorithm where Id=? and State=1`
		err = this.Select(tx, &res, StrSql, id)
		if err != nil {
			return nil, err
		}
	} else {
		StrSql := `select * from tbl_algorithm where State=1`
		err = this.Select(tx, &res, StrSql)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (this *SagaApiDB) QueryEnvsById(tx *sqlx.Tx, id uint32) ([]*tables.Env, error) {
	var err error
	res := make([]*tables.Env, 0)
	if id != 0 {
		StrSql := `select * from tbl_env where Id=? and State=1`
		err = this.Select(tx, &res, StrSql, id)
		if err != nil {
			return nil, err
		}
	} else {
		StrSql := `select * from tbl_env where State=1`
		err = this.Select(tx, &res, StrSql)
		if err != nil {
			return nil, err
		}
	}

	log.Debugf("env id : %d, len(res): %d", id, len(res))
	return res, nil
}
