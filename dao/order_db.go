package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

func (this *SagaApiDB) ClearOrderDB() error {
	strSql := "delete from tbl_order"
	_, err := this.DB.Exec(strSql)
	return err
}

func (this *SagaApiDB) InsertOrder(tx *sqlx.Tx, order *tables.Order) error {
	// use NameExec better.
	strSql := `insert into tbl_order (OrderId,Title, ProductName, OrderType, OrderTime, OrderStatus,Amount, 
OntId,UserName,Price,ApiId,ApiUrl,SpecificationsId,Coin) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, order.OrderId, order.Title, order.ProductName, order.OrderType, order.OrderTime, order.OrderStatus,
		order.Amount, order.OntId, order.UserName, order.Price, order.ApiId, order.ApiUrl, order.SpecificationsId, order.Coin)
	return err
}

func (this *SagaApiDB) UpdateTxInfoByOrderId(tx *sqlx.Tx, orderId string, txHash string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set TxHash=?,OrderStatus=? where OrderId=?"
	err := this.Exec(tx, strSql, txHash, status, orderId)
	return err
}

func (this *SagaApiDB) QueryOrderStatusByOrderId(tx *sqlx.Tx, orderId string) (sagaconfig.OrderStatus, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=?`
	var orderStatus uint8
	err := this.Get(tx, &orderStatus, strSql, orderId)
	if err != nil {
		return 0, err
	}
	return sagaconfig.OrderStatus(orderStatus), nil
}

func (this *SagaApiDB) QueryOrderByOrderId(tx *sqlx.Tx, orderId string) (*tables.Order, error) {
	strSql := `select * from tbl_order where OrderId=?`
	order := &tables.Order{}
	err := this.Get(tx, order, strSql, orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (this *SagaApiDB) QueryOrderSum(tx *sqlx.Tx, ontId string) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=?`
	var sum int
	err := this.Select(tx, &sum, strSql, ontId)
	if err != nil {
		return 0, nil
	}
	return sum, nil
}

func (this *SagaApiDB) QueryOrderByPage(tx *sqlx.Tx, start, pageSize int, ontId string) ([]*tables.Order, error) {
	strSql := `select * from tbl_order where OntId=? order by OrderTime desc limit ?, ?`
	var res []*tables.Order
	err := this.Select(tx, &res, strSql, ontId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *SagaApiDB) UpdateOrderStatus(tx *sqlx.Tx, orderId string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set OrderStatus=? where OrderId=?"
	err := this.Exec(tx, strSql, status, orderId)
	return err
}

func (this *SagaApiDB) DeleteOrderByOrderId(tx *sqlx.Tx, orderId string) error {
	strSql := `delete from tbl_order where OrderId=?`
	err := this.Exec(tx, strSql, orderId)
	return err
}
