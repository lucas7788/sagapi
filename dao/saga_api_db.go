package dao

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/sagaconfig"
	"time"
)

var DefSagaApiDB *SagaApiDB

type SagaApiDB struct {
	ApiDB    *ApiDB
	OrderDB  *OrderDB
	QrCodeDB *QrCodeDB
	OtherDB  *OtherDB
}

func NewSagaApiDB(dbConfig *sagaconfig.DBConfig) (*SagaApiDB, error) {
	dbx, dberr := sqlx.Open("mysql",
		dbConfig.ProjectDBUser+
			":"+dbConfig.ProjectDBPassword+
			"@tcp("+dbConfig.ProjectDBUrl+
			")/"+dbConfig.ProjectDBName+
			"?charset=utf8&parseTime=true&loc=Asia%2FShanghai")
	if dberr != nil {
		return nil, dberr
	}

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	err := dbx.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &SagaApiDB{
		ApiDB:    NewApiDB(dbx),
		OrderDB:  NewOrderDB(dbx),
		QrCodeDB: NewQrCodeDB(dbx),
	}, nil
}

func (this *SagaApiDB) Close() error {
	return this.ApiDB.conn.Close()
}
