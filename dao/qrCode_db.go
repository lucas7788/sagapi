package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
)

func (this *SagaApiDB) InsertQrCode(tx *sqlx.Tx, code *tables.QrCode) error {
	strSql := `insert into tbl_qr_code (QrCodeId,Ver, OrderId, Requester, Signature,Signer,QrCodeData,Callback,Exp,
Chain,QrCodeDesc,ContractType) values (?,?,?,?,?,?,?,?,?,?,?,?)`

	err := this.Exec(tx, strSql, code.QrCodeId, code.Ver, code.OrderId, code.Requester, code.Signature, code.Signer,
		code.QrCodeData, code.Callback, code.Exp, code.Chain, code.QrCodeDesc, code.ContractType)
	return err
}

func (this *SagaApiDB) DeleteQrCodeByOrderId(tx *sqlx.Tx, orderId string) error {
	strSql := `delete from tbl_qr_code where OrderId=?`
	err := this.Exec(tx, strSql, orderId)
	return err
}

func (this *SagaApiDB) QueryOrderIdByQrCodeId(tx *sqlx.Tx, qrCodeId string) (string, error) {
	code, err := this.QueryQrCodeByQrCodeId(tx, qrCodeId)
	if err != nil {
		return "", err
	}
	return code.OrderId, nil
}

func (this *SagaApiDB) QueryQrCodeByOrderId(tx *sqlx.Tx, orderId string) (*tables.QrCode, error) {
	return this.queryQrCodeById(tx, orderId, "")
}

func (this *SagaApiDB) QueryQrCodeByQrCodeId(tx *sqlx.Tx, qrCodeId string) (*tables.QrCode, error) {
	return this.queryQrCodeById(tx, "", qrCodeId)
}

func (this *SagaApiDB) QueryQrCodeResultByQrCodeId(tx *sqlx.Tx, qrCodeId string) (string, error) {
	strSql := `select OrderStatus from tbl_order where OrderId=(select OrderId from tbl_qr_code where QrCodeId=?)`
	var orderStatus uint8
	err := this.Get(tx, &orderStatus, strSql, qrCodeId)
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

func (this *SagaApiDB) queryQrCodeById(tx *sqlx.Tx, orderId, qrCodeId string) (*tables.QrCode, error) {
	var strSql string
	if orderId != "" {
		strSql = `select QrCodeId,Ver,OrderId,Requester,Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc,ContractType  from tbl_qr_code where OrderId=?`
	} else if qrCodeId != "" {
		strSql = `select QrCodeId,Ver,OrderId,Requester,Signature,Signer,QrCodeData,Callback,Exp,Chain,QrCodeDesc,ContractType  from tbl_qr_code where QrCodeId=?`
	}

	var where string
	if orderId != "" {
		where = orderId
	} else if qrCodeId != "" {
		where = qrCodeId
	}
	qrCode := &tables.QrCode{}
	err := this.Get(tx, qrCode, strSql, where)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}
