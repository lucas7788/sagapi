package dao

import (
	"testing"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/stretchr/testify/assert"
	"time"
)

var TestDB *SagaApiDB

func TestMain(m *testing.M) {
	fmt.Println("begin test db.")
	var err error
	TestDB, err = NewSagaApiDB(sagaconfig.DefDBConfigMap[sagaconfig.NETWORK_ID_SOLO_NET])
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Println("end")
}

func TestOrderDB_UpdateOrderStatus(t *testing.T) {
	err := TestDB.OrderDB.UpdateOrderStatus("cc5ebaf6-24cb-423b-a590-e05a41f8c1f5", sagaconfig.Canceled)
	assert.Nil(t, err)
}

func TestOrderDB_InsertOrder(t *testing.T) {
	tt := time.Now().Unix()
	orderId := "abcedkfy"
	order := &tables.Order{
		ApiId:     1,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err := TestDB.OrderDB.InsertOrder(order)
	assert.Nil(t, err)
	orderFromDb, err := TestDB.OrderDB.QueryOrderByOrderId(orderId)
	assert.Nil(t, err)
	assert.Equal(t, orderFromDb.OrderId, orderId)

	orderId2 := "abcedkfyfgt"
	order2 := &tables.Order{
		ApiId:     1,
		OrderId:   orderId2,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.OrderDB.InsertOrder(order2)
	assert.Nil(t, err)

	code := &tables.QrCode{
		QrCodeId: "qbcdab",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.QrCodeDB.InsertQrCode(code)
	assert.Nil(t, err)
	code = &tables.QrCode{
		QrCodeId: "qbcdabc",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.QrCodeDB.InsertQrCode(code)
	assert.Nil(t, err)
	code, err = TestDB.QrCodeDB.QueryQrCodeByOrderId(orderId)
	assert.Nil(t, err)
	fmt.Println(code)
	err = TestDB.QrCodeDB.DeleteQrCodeByOrderId(orderId)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId)
	assert.Nil(t, err)
}

func TestOrderDB_QueryOrderByOrderId(t *testing.T) {
	orderId := "abcedkfy"
	err := TestDB.OrderDB.UpdateTxInfoByOrderId(orderId, "", sagaconfig.Canceled)
	assert.Nil(t, err)

	order, err := TestDB.OrderDB.QueryOrderByOrderId(orderId)
	assert.Nil(t, err)
	assert.Equal(t, order.OrderId, orderId)
	assert.Equal(t, order.OrderStatus, sagaconfig.Canceled)

	orders, err := TestDB.OrderDB.QueryOrderByPage(0, 1, "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(orders))

	err = TestDB.OrderDB.UpdateOrderStatus(orderId, sagaconfig.Completed)
	assert.Nil(t, err)
	status, err := TestDB.OrderDB.QueryOrderStatusByOrderId(orderId)
	assert.Nil(t, err)
	assert.Equal(t, sagaconfig.Completed, status)
}

func TestOrderDB_DeleteOrderByOrderId(t *testing.T) {
	orderId := "abcedkfyfgtghj"
	tt := time.Now().Unix()
	order2 := &tables.Order{
		ApiId:     1,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err := TestDB.OrderDB.InsertOrder(order2)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId)
	assert.Nil(t, err)
}

func TestSagaDB_QueryOrderStatusByOrderId(t *testing.T) {
	_, err := TestDB.OrderDB.QueryOrderByOrderId("1")
	assert.NotNil(t, err)
}

func TestOrderDB_QueryOrderByPage(t *testing.T) {
	order, err := TestDB.OrderDB.QueryOrderByPage(1, 2, "")
	assert.Nil(t, err)
	fmt.Println(order)
}

func TestOrderDB_UpdateTxInfoByOrderId(t *testing.T) {
	err := TestDB.OrderDB.UpdateTxInfoByOrderId("145e89f6-850e-44a7-be3e-9224fd066858", "", sagaconfig.Processing)
	assert.Nil(t, err)
}

func TestOrderDB_QueryQrCodeResultByQrCodeId(t *testing.T) {
	status, err := TestDB.QrCodeDB.QueryQrCodeResultByQrCodeId("bb5b68d4-0282-469d-936b-ae43e30c5de5")
	assert.Nil(t, err)
	fmt.Println(status)
}
