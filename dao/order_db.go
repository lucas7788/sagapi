package dao

import (
	"database/sql"
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

func (this *OrderDB) InsertOrder(buyRecord *tables.Order) error {
	strSql := `insert into order_tbl (OrderId, ProductName, OrderType, OrderTime, PayTime, OrderStatus,Amount, 
OntId,UserName,TxHash,Price,ApiId,ApiKey,Specifications) values (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(buyRecord.OrderId, buyRecord.ProductName, buyRecord.OrderType, buyRecord.OrderTime, buyRecord.PayTime,
		buyRecord.OrderStatus, buyRecord.Amount, buyRecord.OntId, buyRecord.UserName, buyRecord.TxHash, buyRecord.Price, buyRecord.ApiId,
		buyRecord.ApiKey, buyRecord.Specifications)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrderDB) UpdateTxInfoByOrderId(orderId string, txHash string, status config.OrderStatus) error {
	strSql := "update order_tbl set TxHash=?,Status=? where OrderId=?"
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

func (this *OrderDB) QueryOrderByOrderId(orderId string) (*tables.Order, error) {
	strSql := `select OrderId, ProductName, OrderType, OrderTime, PayTime, OrderStatus,Amount, 
OntId,UserName,TxHash,Price,ApiId,ApiKey,Specifications from order_tbl where OrderId=?`
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
		var orderTime, paiedTime int64
		var orderId, productName, orderType, amount, ontId, userName, txHash, price, apiKey string
		var specifications, apiId int
		var orderStatus uint8
		if err = rows.Scan(&orderId, &productName, &orderType, &orderTime, &paiedTime, &orderStatus, &amount,
			&ontId, &userName, &txHash, &price, &apiId, &apiKey, &specifications); err != nil {
			return nil, err
		}
		return &tables.Order{
			OrderId:        orderId,
			ProductName:    productName,
			OrderType:      orderType,
			OrderTime:      orderTime,
			PayTime:        paiedTime,
			OrderStatus:    config.OrderStatus(orderStatus),
			Amount:         amount,
			OntId:          ontId,
			UserName:       userName,
			TxHash:         txHash,
			Price:          price,
			ApiId:          apiId,
			ApiKey:         apiKey,
			Specifications: specifications,
		}, nil
	}
	return nil, nil
}

func (this *OrderDB) InsertQrCode(code *tables.QrCode) error {

	strSql := `insert into qr_code_tbl (Ver, Id, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc) values (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := this.db.Prepare(strSql)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(code.Ver, code.QrCodeId, code.OrderId, code.Requester, code.Signature, code.Signer, code.QrCodeData, code.Callback, code.Exp, code.Chain, code.QrCodeDesc)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrderDB) QueryOrderIdByQrCodeId(qrCodeId string) (string, error) {
	code, err := this.QueryQrCodeByQrCodeId(qrCodeId)
	if err != nil {
		return "", err
	}
	return code.OrderId, nil
}

func (this *OrderDB) UpdateOrderStatus(orderId string, status config.OrderStatus) error {
	strSql := "update order_tbl set Status=? where OrderId=?"
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

func (this *OrderDB) DeleteOrderByOrderId(orderId string) error {
	strSql := `delete from qr_code_tbl where OrderId=?`
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
		strSql = `select Ver, Id, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from qr_code_tbl where OrderId=?`
	} else if qrCodeId != "" {
		strSql = `select Ver, Id, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from qr_code_tbl where QrCodeId=?`
	}

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
		var exp int64
		var ver, id, orderId, requester, signature, signer, qrCodeData, callback, chain, qrCodeDesc string
		if err = rows.Scan(&ver, &id, &orderId, &requester, &signature, &signer, &qrCodeData, &callback, &exp, &chain, &qrCodeDesc); err != nil {
			return nil, err
		}
		return &tables.QrCode{
			Ver:        ver,
			QrCodeId:         id,
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
	return nil, nil
}
func (this *OrderDB) Close() error {
	return this.db.Close()
}
