package dao

import (
	"context"
	"database/sql"
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
	db, dberr := sql.Open("mysql",
		dbConfig.ProjectDBUser+
			":"+dbConfig.ProjectDBPassword+
			"@tcp("+dbConfig.ProjectDBUrl+
			")/"+dbConfig.ProjectDBName+
			"?charset=utf8&parseTime=true")
	if dberr != nil {
		return nil, dberr
	}
	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &SagaApiDB{
		ApiDB:    NewApiDB(db),
		OrderDB:  NewOrderDB(db),
		QrCodeDB: NewQrCodeDB(db),
		OtherDB:  NewOtherDB(db),
	}, nil
}

func (this *SagaApiDB) Close() error {
	return this.ApiDB.db.Close()
}
