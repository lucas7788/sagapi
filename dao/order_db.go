package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

func (this *SagaApiDB) ClearOrderDB() error {
	strSql := "delete from tbl_order"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *SagaApiDB) InsertOrder(order *tables.Order) error {
	strSql := `insert into tbl_order (OrderId,Title, ProductName, OrderType, OrderTime, OrderStatus,Amount, 
OntId,UserName,Price,ApiId,ApiUrl,SpecificationsId,Coin) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := this.DB.Exec(strSql, order.OrderId, order.Title, order.ProductName, order.OrderType, order.OrderTime, order.OrderStatus,
		order.Amount, order.OntId, order.UserName, order.Price, order.ApiId, order.ApiUrl, order.SpecificationsId, order.Coin)
	return err
}

func (this *SagaApiDB) UpdateTxInfoByOrderId(orderId string, txHash string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set TxHash=?,OrderStatus=? where OrderId=?"
	_, err := this.DB.Exec(strSql, txHash, status, orderId)
	return err
}

func (this *SagaApiDB) QueryOrderStatusByOrderId(orderId string) (sagaconfig.OrderStatus, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=?`
	var orderStatus uint8
	err := this.DB.Get(&orderStatus, strSql, orderId)
	if err != nil {
		return 0, err
	}
	return sagaconfig.OrderStatus(orderStatus), nil
}

func (this *SagaApiDB) QueryOrderByOrderId(orderId string) (*tables.Order, error) {
	strSql := `select * from tbl_order where OrderId=?`
	order := &tables.Order{}
	err := this.DB.Get(order, strSql, orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (this *SagaApiDB) QueryOrderSum(ontId string) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=?`
	var sum int
	err := this.DB.Select(&sum, strSql, ontId)
	if err != nil {
		return 0, nil
	}
	return sum, nil
}

func (this *SagaApiDB) QueryOrderByPage(start, pageSize int, ontId string) ([]*tables.Order, error) {
	strSql := `select * from tbl_order where OntId=? order by OrderTime desc limit ?, ?`
	var res []*tables.Order
	err := this.DB.Select(&res, strSql, ontId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *SagaApiDB) UpdateOrderStatus(orderId string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set OrderStatus=? where OrderId=?"
	_, err := this.DB.Exec(strSql, status, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (this *SagaApiDB) DeleteOrderByOrderId(orderId string) error {
	strSql := `delete from tbl_order where OrderId=?`
	_, err := this.DB.Exec(strSql, orderId)
	if err != nil {
		return err
	}
	return nil
}
