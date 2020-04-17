package dao

import (
	"database/sql"
	"github.com/ontio/sagapi/config"
)

var DefSagaApiDB *SagaApiDB

type SagaApiDB struct {
	ApiDB   *ApiDB
	OrderDB *OrderDB
}

func NewSagaApiDB() (*SagaApiDB, error) {
	db, dberr := sql.Open("mysql",
		config.DefConfig.DbConfig.ProjectDBUser+
			":"+config.DefConfig.DbConfig.ProjectDBPassword+
			"@tcp("+config.DefConfig.DbConfig.ProjectDBUrl+
			")/"+config.DefConfig.DbConfig.ProjectDBName+
			"?charset=utf8")
	if dberr != nil {
		return nil, dberr
	}
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return &SagaApiDB{
		ApiDB:   NewApiDB(db),
		OrderDB: NewOrderDB(db),
	}, nil
}

func (this *SagaApiDB) Close() error {
	return this.ApiDB.db.Close()
}
