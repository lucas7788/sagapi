package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/sagaconfig"
)

type OtherDB struct {
	db *sql.DB
}

func NewOtherDB(db *sql.DB) *OtherDB {
	return &OtherDB{
		db: db,
	}
}

func (this *OtherDB) TblApiKeyOfUpdateOrderStatus(orderId string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_api_key set OrderStatus=? where OrderId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, orderId)
	return err
}
