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
	TestDB, err = NewSagaApiDB(sagaconfig.DefDBConfigMap[sagaconfig.NETWORK_ID_POLARIS_NET])
	if err != nil {
		panic(err)
	}
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
		ApiProvider:     "",
		ApiUrl:          "",
		Price:           "",
		ApiDesc:         "",
		Specifications:  1,
		Popularity:      0,
		Delay:           0,
		SuccessRate:     0,
		InvokeFrequency: 0,
	}
	err = TestDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{info})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)
	detail := &tables.ApiDetailInfo{
		ApiId:               basic[0].ApiId,
		RequestType:         "POST",
		Mark:                "",
		ResponseParam:       "",
		ResponseExample:     "",
		DataDesc:            "",
		DataSource:          "",
		ApplicationScenario: "",
	}
	TestDB.ApiDB.InsertApiDetailInfo(detail)
	m.Run()
	fmt.Println("end")
	err = TestDB.ApiDB.ClearApiDetailDB()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	err = TestDB.ApiDB.ClearApiBasicDB()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}

func TestOrderDB_InsertOrder(t *testing.T) {
	tt := time.Now().Unix()
	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)
	assert.Nil(t, err)
	orderId := "abcedkfy"
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.OrderDB.InsertOrder(order)
	assert.Nil(t, err)
	orderFromDb, err := TestDB.OrderDB.QueryOrderByOrderId(orderId)
	assert.Nil(t, err)
	assert.Equal(t, orderFromDb.OrderId, orderId)

	err = TestDB.OrderDB.UpdateOrderStatus(orderId, sagaconfig.Canceled)
	assert.Nil(t, err)

	orderId2 := "abcedkfyfgt"
	order2 := &tables.Order{
		ApiId:     basic[0].ApiId,
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
	_, err = TestDB.QrCodeDB.QueryQrCodeResultByQrCodeId(code.QrCodeId)
	assert.Nil(t, err)

	err = TestDB.QrCodeDB.DeleteQrCodeByOrderId(orderId)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId2)
	assert.Nil(t, err)
}

func TestOrderDB_QueryOrderByOrderId(t *testing.T) {
	TestDB.OrderDB.ClearOrderDB()
	tt := time.Now().Unix()
	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)
	assert.Nil(t, err)
	orderId := "abcedkfy"
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.OrderDB.InsertOrder(order)
	assert.Nil(t, err)
	err = TestDB.OrderDB.UpdateTxInfoByOrderId(orderId, "", sagaconfig.Canceled)
	assert.Nil(t, err)

	order, err = TestDB.OrderDB.QueryOrderByOrderId(orderId)
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

	TestDB.OrderDB.ClearOrderDB()
}

func TestOrderDB_DeleteOrderByOrderId(t *testing.T) {
	orderId := "abcedkfyfgtghj"
	basic, err := TestDB.ApiDB.QueryApiBasicInfoByPage(0, 1)
	assert.Nil(t, err)
	tt := time.Now().Unix()
	order2 := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.OrderDB.InsertOrder(order2)
	assert.Nil(t, err)
	err = TestDB.OrderDB.DeleteOrderByOrderId(orderId)
	assert.Nil(t, err)
}

func TestOrderDB_QueryOrderByPage(t *testing.T) {
	order, err := TestDB.OrderDB.QueryOrderByPage(0, 2, "")
	assert.Nil(t, err)
	fmt.Println(order)
}
