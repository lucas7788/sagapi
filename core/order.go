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

func (this *SagaOrder) TakeOrder(param *common.TakeOrderParam) (*tables.QrCode, error) {
	info, err := dao.DefDB.QueryApiBasicInfoByApiId(uint(param.ApiId))
	if err != nil {
		return nil, err
	}
	price := utils.ToIntByPrecise(info.ApiPrice, config.ONG_DECIMALS)
	specifications := new(big.Int).SetUint64(uint64(param.Specifications))
	amount := new(big.Int).Mul(price, specifications)
	amountStr := utils.ToStringByPrecise(amount, config.ONG_DECIMALS)
	orderId := common.GenerateOrderId()
	order := &tables.Order{
		OrderId:     orderId,
		ProductName: param.ProductName,
		Type:        config.ApiOrder,
		OrderTime:   time.Now().Unix(),
		OrderStatus: config.Processing,
		Amount:      amountStr,
		OntId:       param.OntId,
		UserName:    param.UserName,
		Price:       info.ApiPrice,
		ApiId:       info.ApiId,
	}
	err = dao.DefDB.InsertOrder(order)
	if err != nil {
		return nil, err
	}
	code := common.BuildTestNetQrCode(orderId, "", "", param.OntId, "", amountStr)
	err = dao.DefDB.InsertQrCode(code)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func (this *SagaOrder) GetPayQrCodeById(id string) (*tables.QrCode, error) {
	return dao.DefDB.QueryQrCodeById(id)
}
