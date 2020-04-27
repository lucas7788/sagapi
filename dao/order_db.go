package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

type OrderDB struct {
	conn *sqlx.DB
}

func NewOrderDB(cn *sqlx.DB) *OrderDB {
	return &OrderDB{
		conn: cn,
	}
}
