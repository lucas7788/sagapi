package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

type QrCodeDB struct {
	db *sqlx.DB
}

func NewQrCodeDB(db *sqlx.DB) *QrCodeDB {
	return &QrCodeDB{
		db: db,
	}
}

func (this *QrCodeDB) InsertQrCode(code *tables.QrCode) error {
	strSql := `insert into tbl_qr_code (QrCodeId,Ver, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,
Chain,QrCodeDesc) values (?,?,?,?,?,?,?,?,?,?,?)`

	_, err := this.db.Exec(strSql, code.QrCodeId, code.Ver, code.OrderId, code.Requester, code.Signature, code.Signer,
		code.QrCodeData, code.Callback, code.Exp, code.Chain, code.QrCodeDesc)
	if err != nil {
		return err
	}
	return nil
}

func (this *QrCodeDB) DeleteQrCodeByOrderId(orderId string) error {
	strSql := `delete from tbl_qr_code where OrderId=?`
	_, err := this.db.Exec(strSql, orderId)
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
	var orderStatus uint8
	err := this.db.Get(&orderStatus, strSql, qrCodeId)
	if err != nil {
		return "", err
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

func (this *QrCodeDB) queryQrCodeById(orderId, qrCodeId string) (*tables.QrCode, error) {
	var strSql string
	if orderId != "" {
		strSql = `select QrCodeId,Ver,OrderId,Requester,Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from tbl_qr_code where OrderId=?`
	} else if qrCodeId != "" {
		strSql = `select QrCodeId,Ver,OrderId,Requester,Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc from tbl_qr_code where QrCodeId=?`
	}

	var where string
	if orderId != "" {
		where = orderId
	} else if qrCodeId != "" {
		where = qrCodeId
	}
	qrCode := &tables.QrCode{}
	err := this.db.Get(qrCode, strSql, where)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

func (this *QrCodeDB) Close() error {
	return this.db.Close()
}
