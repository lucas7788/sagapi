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

var DefSagaApi *SagaApi

type SagaApi struct {
}

func (this *SagaApi) TakeOrder(param *common.TakeOrderParam) error {
	info, err := dao.DefDB.QueryApiBasicInfoByApiId(uint(param.ApiId))
	if err != nil {
		return err
	}
	price := utils.ToIntByPrecise(info.ApiPrice, config.ONG_DECIMALS)
	specifications := new(big.Int).SetUint64(uint64(param.Specifications))
	amount := new(big.Int).Mul(price, specifications)
	orderId := common.GenerateOrderId()
	order := &tables.Order{
		OrderId:     orderId,
		ProductName: param.ProductName,
		Type:        config.ApiOrder,
		OrderTime:   time.Now().Unix(),
		OrderStatus: config.Processing,
		Amount:      utils.ToStringByPrecise(amount, config.ONG_DECIMALS),
		OntId:       param.OntId,
		UserName:    param.UserName,
		Price:       info.ApiPrice,
		ApiId:       info.ApiId,
		ApiKey:      "",
	}
	return dao.DefDB.InsertOrder(order)
}
