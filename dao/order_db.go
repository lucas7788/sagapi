package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

type OrderDB struct {
	db *sql.DB
}

func NewOrderDB(db *sql.DB) *OrderDB {
	return &OrderDB{
		db: db,
	}
}

func (this *OrderDB) InsertOrder(order *tables.Order) error {
	strSql := `insert into tbl_order (OrderId,Title, ProductName, OrderType, OrderTime, OrderStatus,Amount, 
OntId,UserName,Price,ApiId,SpecificationsId,Coin) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.OrderId, order.Title, order.ProductName, order.OrderType, order.OrderTime, order.OrderStatus,
		order.Amount, order.OntId, order.UserName, order.Price, order.ApiId, order.SpecificationsId, order.Coin)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrderDB) UpdateTxInfoByOrderId(orderId string, txHash string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set TxHash=?,OrderStatus=? where OrderId=?"
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(txHash, status, orderId)
	return err
}

func (this *OrderDB) QueryOrderStatusByOrderId(orderId string) (sagaconfig.OrderStatus, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query(orderId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var orderStatus uint8
		if err = rows.Scan(&orderStatus); err != nil {
			return 0, err
		}
		return sagaconfig.OrderStatus(orderStatus), nil
	}
	return 0, fmt.Errorf("order not found, orderId: %s", orderId)
}

func (this *OrderDB) QueryOrderByOrderId(orderId string) (*tables.Order, error) {
	strSql := `select OrderId,Title, ProductName, OrderType, OrderTime, PayTime, OrderStatus,Amount, 
OntId,UserName,TxHash,Price,ApiId,SpecificationsId,Coin from tbl_order where OrderId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(orderId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var orderTime, payTime int64
		var orderId, title, productName, orderType, amount, ontId, userName, txHash, price, coin string
		var specifications, apiId int
		var orderStatus uint8
		if err = rows.Scan(&orderId, &title, &productName, &orderType, &orderTime, &payTime, &orderStatus, &amount,
			&ontId, &userName, &txHash, &price, &apiId, &specifications, &coin); err != nil {
			return nil, err
		}
		return &tables.Order{
			OrderId:          orderId,
			Title:            title,
			ProductName:      productName,
			OrderType:        orderType,
			OrderTime:        orderTime,
			PayTime:          payTime,
			OrderStatus:      sagaconfig.OrderStatus(orderStatus),
			Amount:           amount,
			OntId:            ontId,
			UserName:         userName,
			TxHash:           txHash,
			Price:            price,
			ApiId:            apiId,
			SpecificationsId: specifications,
		}, nil
	}
	return nil, fmt.Errorf("order not found, orderId: %s", orderId)
}

func (this *OrderDB) QueryOrderSum(ontId string) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query(ontId)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var sum int
		if err := rows.Scan(&sum); err != nil {
			return 0, err
		}
		return sum, nil
	}
	return 0, nil
}
func (this *OrderDB) QueryOrderByPage(start, pageSize int, ontId string) ([]*tables.Order, error) {
	strSql := `select OrderId,Title, ProductName, OrderType, OrderTime, PayTime, OrderStatus,Amount, 
OntId,UserName,TxHash,Price,ApiId,SpecificationsId,Coin from tbl_order where OntId=? order by OrderTime desc limit ?, ?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(ontId, start, pageSize)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	res := make([]*tables.Order, 0)
	for rows.Next() {
		var orderTime, payTime int64
		var orderId, title, productName, orderType, amount, ontId, userName, txHash, price, coin string
		var specifications, apiId int
		var orderStatus uint8
		if err = rows.Scan(&orderId, &title, &productName, &orderType, &orderTime, &payTime, &orderStatus, &amount,
			&ontId, &userName, &txHash, &price, &apiId, &specifications, &coin); err != nil {
			return nil, err
		}
		res = append(res, &tables.Order{
			OrderId:          orderId,
			Title:            title,
			ProductName:      productName,
			OrderType:        orderType,
			OrderTime:        orderTime,
			PayTime:          payTime,
			OrderStatus:      sagaconfig.OrderStatus(orderStatus),
			Amount:           amount,
			OntId:            ontId,
			UserName:         userName,
			TxHash:           txHash,
			Price:            price,
			ApiId:            apiId,
			SpecificationsId: specifications,
			Coin:             coin,
		})
	}
	return res, nil
}

func (this *OrderDB) UpdateOrderStatus(orderId string, status sagaconfig.OrderStatus) error {
	strSql := "update tbl_order set OrderStatus=? where OrderId=?"
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

func (this *OrderDB) DeleteOrderByOrderId(orderId string) error {
	strSql := `delete from tbl_order where OrderId=?`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(orderId)
	return err
}
