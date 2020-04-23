package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
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

func (this *OrderDB) UpdateTxInfoByOrderId(orderId string, txHash string, status config.OrderStatus) error {
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

func (this *OrderDB) QueryOrderStatusByOrderId(orderId string) (config.OrderStatus, error) {
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
		return config.OrderStatus(orderStatus), nil
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
			OrderStatus:      config.OrderStatus(orderStatus),
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
			OrderStatus:      config.OrderStatus(orderStatus),
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

func (this *OrderDB) InsertQrCode(code *tables.QrCode) error {
	strSql := `insert into tbl_qr_code (QrCodeId,Ver, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,
Chain,QrCodeDesc) values (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(code.QrCodeId, code.Ver, code.OrderId, code.Requester, code.Signature, code.Signer,
		code.QrCodeData, code.Callback, code.Exp, code.Chain, code.QrCodeDesc)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrderDB) DeleteQrCodeByOrderId(orderId string) error {
	strSql := `delete from tbl_qr_code where OrderId=?`
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

func (this *OrderDB) QueryOrderIdByQrCodeId(qrCodeId string) (string, error) {
	code, err := this.QueryQrCodeByQrCodeId(qrCodeId)
	if err != nil {
		return "", err
	}
	return code.OrderId, nil
}

func (this *OrderDB) UpdateOrderStatus(orderId string, status config.OrderStatus) error {
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
func (this *OrderDB) UpdateOrderStatusInApiKey(orderId string, status config.OrderStatus) error {
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

func (this *OrderDB) QueryQrCodeByOrderId(orderId string) (*tables.QrCode, error) {
	return this.queryQrCodeById(orderId, "")
}

func (this *OrderDB) QueryQrCodeByQrCodeId(qrCodeId string) (*tables.QrCode, error) {
	return this.queryQrCodeById("", qrCodeId)
}

func (this *OrderDB) QueryQrCodeResultByQrCodeId(qrCodeId string) (string, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=(select OrderId from tbl_qr_code where QrCodeId=?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return "", err
	}
	rows, err := stmt.Query(qrCodeId)
	if err != nil {
		return "", err
	}
	if rows != nil {
		defer rows.Close()
	}
	for rows.Next() {
		var orderStatus uint8
		if err = rows.Scan(&orderStatus); err != nil {
			return "0", err //processing
		}
		if orderStatus == uint8(config.Completed) {
			return "1", nil //success
		}
		if orderStatus == uint8(config.Processing) {
			return "", nil
		}
		if orderStatus == uint8(config.Canceled) {
			return "", nil //failed
		}
		if orderStatus == uint8(config.Failed) {
			return "0", nil
		}
		return "", nil
	}
	return "", nil
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
func (this *OrderDB) queryQrCodeById(orderId, qrCodeId string) (*tables.QrCode, error) {
	var strSql string
	if orderId != "" {
		strSql = `select Ver, QrCodeId, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from tbl_qr_code where OrderId=?`
	} else if qrCodeId != "" {
		strSql = `select Ver, QrCodeId, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from tbl_qr_code where QrCodeId=?`
	}

	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if orderId != "" {
		rows, err = stmt.Query(orderId)
	} else if qrCodeId != "" {
		rows, err = stmt.Query(qrCodeId)
	}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var exp int64
		var ver, id, orderId, requester, signature, signer, qrCodeData, callback, chain, qrCodeDesc string
		if err = rows.Scan(&ver, &id, &orderId, &requester, &signature, &signer, &qrCodeData, &callback, &exp, &chain, &qrCodeDesc); err != nil {
			return nil, err
		}
		return &tables.QrCode{
			Ver:        ver,
			QrCodeId:   id,
			OrderId:    orderId,
			Requester:  requester,
			Signature:  signature,
			Signer:     signer,
			QrCodeData: qrCodeData,
			Callback:   callback,
			Exp:        exp,
			Chain:      chain,
			QrCodeDesc: qrCodeDesc,
		}, nil
	}
	return nil, fmt.Errorf("not found")
}
func (this *OrderDB) Close() error {
	return this.db.Close()
}
