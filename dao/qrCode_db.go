package dao

import (
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

type QrCodeDB struct {
	db *sql.DB
}

func NewQrCodeDB(db *sql.DB) *QrCodeDB {
	return &QrCodeDB{
		db: db,
	}
}

func (this *QrCodeDB) InsertQrCode(code *tables.QrCode) error {
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

func (this *QrCodeDB) DeleteQrCodeByOrderId(orderId string) error {
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

func (this *QrCodeDB) QueryOrderIdByQrCodeId(qrCodeId string) (string, error) {
	code, err := this.QueryQrCodeByQrCodeId(qrCodeId)
	if err != nil {
		return "", err
	}
	return code.OrderId, nil
}

func (this *QrCodeDB) QueryQrCodeByOrderId(orderId string) (*tables.QrCode, error) {
	return this.queryQrCodeById(orderId, "")
}

func (this *QrCodeDB) QueryQrCodeByQrCodeId(qrCodeId string) (*tables.QrCode, error) {
	return this.queryQrCodeById("", qrCodeId)
}

func (this *QrCodeDB) QueryQrCodeResultByQrCodeId(qrCodeId string) (string, error) {
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
		if orderStatus == uint8(sagaconfig.Completed) {
			return "1", nil //success
		}
		if orderStatus == uint8(sagaconfig.Processing) {
			return "", nil
		}
		if orderStatus == uint8(sagaconfig.Canceled) {
			return "", nil //failed
		}
		if orderStatus == uint8(sagaconfig.Failed) {
			return "0", nil
		}
		return "", nil
	}
	return "", nil
}

func (this *QrCodeDB) queryQrCodeById(orderId, qrCodeId string) (*tables.QrCode, error) {
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

func (this *QrCodeDB) Close() error {
	return this.db.Close()
}
