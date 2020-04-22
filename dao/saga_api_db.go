package dao

import (
	"database/sql"
	"github.com/ontio/sagapi/config"
	"context"
	"time"
)

var DefSagaApiDB *SagaApiDB

type SagaApiDB struct {
	ApiDB   *ApiDB
	OrderDB *OrderDB
}

func NewSagaApiDB() (*SagaApiDB, error) {
	db, dberr := sql.Open("mysql",
		config.DefSagaConfig.DbConfig.ProjectDBUser+
			":"+config.DefSagaConfig.DbConfig.ProjectDBPassword+
			"@tcp("+config.DefSagaConfig.DbConfig.ProjectDBUrl+
			")/"+config.DefSagaConfig.DbConfig.ProjectDBName+
			"?charset=utf8")
	if dberr != nil {
		return nil, dberr
	}
	ctx,cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()
	err := db.PingContext(ctx)
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
