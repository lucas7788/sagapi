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
func (this *OrderDB) ClearOrderDB() error {
	strSql := "delete from tbl_order"
	_, err := this.conn.Exec(strSql)
	return err
}
func (this *OrderDB) InsertOrder(order *tables.Order) error {
	strSql := `insert into tbl_order (OrderId,Title, ProductName, OrderType, OrderTime, OrderStatus,Amount, 
OntId,UserName,Price,ApiId,ApiUrl,SpecificationsId,Coin) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := this.conn.Exec(strSql, order.OrderId, order.Title, order.ProductName, order.OrderType, order.OrderTime, order.OrderStatus,
		order.Amount, order.OntId, order.UserName, order.Price, order.ApiId, order.ApiUrl, order.SpecificationsId, order.Coin)
	return err
}

func (this *OrderDB) UpdateTxInfoByOrderId(orderId string, txHash string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set TxHash=?,OrderStatus=? where OrderId=?"
	_, err := this.conn.Exec(strSql, txHash, status, orderId)
	return err
}

func (this *OrderDB) QueryOrderStatusByOrderId(orderId string) (sagaconfig.OrderStatus, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=?`
	var orderStatus uint8
	err := this.conn.Get(&orderStatus, strSql, orderId)
	if err != nil {
		return 0, err
	}
	return sagaconfig.OrderStatus(orderStatus), nil
}

func (this *OrderDB) QueryOrderByOrderId(orderId string) (*tables.Order, error) {
	strSql := `select * from tbl_order where OrderId=?`
	order := &tables.Order{}
	err := this.conn.Get(order, strSql, orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (this *OrderDB) QueryOrderSum(ontId string) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=?`
	var sum int
	err := this.conn.Select(&sum, strSql, ontId)
	if err != nil {
		return 0, nil
	}
	return sum, nil
}

func (this *OrderDB) QueryOrderByPage(start, pageSize int, ontId string) ([]*tables.Order, error) {
	strSql := `select * from tbl_order where OntId=? order by OrderTime desc limit ?, ?`
	var res []*tables.Order
	err := this.conn.Select(&res, strSql, ontId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *OrderDB) UpdateOrderStatus(orderId string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set OrderStatus=? where OrderId=?"
	_, err := this.conn.Exec(strSql, status, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrderDB) DeleteOrderByOrderId(orderId string) error {
	strSql := `delete from tbl_order where OrderId=?`
	_, err := this.conn.Exec(strSql, orderId)
	if err != nil {
		return err
	}
	return nil
}
