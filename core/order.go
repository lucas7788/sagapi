package core

import (
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/utils"
	"math/big"
	"time"
)

var DefSagaOrder *SagaOrder

type SagaOrder struct {
}

func NewSagaOrder() *SagaOrder {
	return &SagaOrder{}
}

func (this *SagaOrder) TakeOrder(param *common.TakeOrderParam) (*common.QrCodeResult, error) {
	info, err := dao.DefDB.QueryApiBasicInfoByApiId(uint(param.ApiId))
	if err != nil {
		return nil, err
	}
	price := utils.ToIntByPrecise(info.ApiPrice, config.ONG_DECIMALS)
	specifications := new(big.Int).SetUint64(uint64(param.Specifications))
	amount := new(big.Int).Mul(price, specifications)
	amountStr := utils.ToStringByPrecise(amount, config.ONG_DECIMALS)
	orderId := common.GenerateUUId()
	order := &tables.Order{
		OrderId:        orderId,
		ProductName:    param.ProductName,
		Type:           config.ApiOrder,
		OrderTime:      time.Now().Unix(),
		OrderStatus:    config.Processing,
		Amount:         amountStr,
		OntId:          param.OntId,
		UserName:       param.UserName,
		Price:          info.ApiPrice,
		ApiId:          info.ApiId,
		Specifications: param.Specifications,
	}
	err = dao.DefDB.InsertOrder(order)
	if err != nil {
		return nil, err
	}
	code := common.BuildTestNetQrCode(orderId, param.OntId, "", param.OntId, "", amountStr)
	err = dao.DefDB.InsertQrCode(code)
	if err != nil {
		return nil, err
	}
	return common.BuildQrCodeResult(code.Id), nil
}

func (this *SagaOrder) GetQrCodeByOrderId(orderId string) (*common.QrCodeResult, error) {
	code, err := dao.DefDB.QueryQrCodeByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	return common.BuildQrCodeResult(code.Id), nil
}

func (this *SagaOrder) GetQrCodeDataById(id string) (*tables.QrCode, error) {
	return dao.DefDB.QueryQrCodeByQrCodeId(id)
}

func (this *SagaOrder) CancelOrder(param *common.OrderIdParam) error {
	return dao.DefDB.UpdateOrderStatus(param.OrderId, config.Canceled)
}
func (this *SagaOrder) DeleteOrderByOrderId(param *common.OrderIdParam) error {
	return dao.DefDB.DeleteOrderByOrderId(param.OrderId)
}

func (this *SagaOrder) GetTxResult(orderId string) (*common.GetOrderResult, error) {
	order, err := dao.DefDB.QueryOrderByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	res := &common.GetOrderResult{
		UserName: order.UserName,
		OntId:    order.OntId,
	}
	if order.OrderStatus == config.Completed {
		res.Result = "1"
	} else if order.OrderStatus == config.Processing {
		res.Result = ""
	} else if order.OrderStatus == config.Failed {
		res.Result = "0"
	}
	return res, nil
}
