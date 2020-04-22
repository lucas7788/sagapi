package core

import (
	"fmt"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/utils"
	"math/big"
	"strings"
	"time"
)

type SagaOrder struct {
}

func NewSagaOrder() *SagaOrder {
	return &SagaOrder{}
}

func (this *SagaOrder) TakeOrder(param *common.TakeOrderParam) (*common.QrCodeResponse, error) {
	info, err := dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoByApiId(param.ApiId)
	if err != nil {
		return nil, err
	}
	spec, err := dao.DefSagaApiDB.ApiDB.QuerySpecificationsBySpecificationsId(param.SpecificationsId)
	if err != nil {
		return nil, err
	}
	price := utils.ToIntByPrecise(spec.Price, config.ONG_DECIMALS)
	specifications := new(big.Int).SetUint64(uint64(spec.Amount))
	amount := new(big.Int).Mul(price, specifications)
	amountStr := utils.ToStringByPrecise(amount, config.ONG_DECIMALS)
	orderId := common.GenerateUUId()
	order := &tables.Order{
		OrderId:          orderId,
		Title:            info.Title,
		ProductName:      param.ProductName,
		OrderType:        config.Api,
		OrderTime:        time.Now().Unix(),
		OrderStatus:      config.Processing,
		Amount:           amountStr,
		OntId:            param.OntId,
		UserName:         param.UserName,
		Price:            spec.Price,
		ApiId:            info.ApiId,
		SpecificationsId: param.SpecificationsId,
		Coin:             config.TOKEN_TYPE_ONG,
	}
	err = dao.DefSagaApiDB.OrderDB.InsertOrder(order)
	if err != nil {
		return nil, err
	}
	arr := strings.Split(param.OntId, ":")
	if len(arr) < 3 {
		return nil, fmt.Errorf("error ontid: %s", param.OntId)
	}
	code := common.BuildTestNetQrCode(orderId, param.OntId, arr[2], arr[2], "AbtTQJYKfQxq4UdygDsbLVjE8uRrJ2H3tP", amountStr)
	err = dao.DefSagaApiDB.OrderDB.InsertQrCode(code)
	if err != nil {
		return nil, err
	}
	return common.BuildQrCodeResult(code.QrCodeId), nil
}

func (this *SagaOrder) QueryOrderByPage(pageNum, pageSize int, ontid string) (map[string]interface{}, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 0 {
		pageSize = 0
	}
	total, err := dao.DefSagaApiDB.OrderDB.QueryOrderSum(ontid)
	if err != nil {
		return nil, err
	}
	start := (pageNum - 1) * pageSize
	orders, err := dao.DefSagaApiDB.OrderDB.QueryOrderByPage(start, pageSize, ontid)
	if err != nil {
		return nil, err
	}
	res := make([]*common.OrderResult, len(orders))
	for i, order := range orders {
		apiKey, err := dao.DefSagaApiDB.ApiDB.QueryApiKeyByOrderId(order.OrderId)
		if err != nil && !strings.Contains(err.Error(),"not found") {
			return nil, err
		}
		if apiKey == nil {
			spec, err := dao.DefSagaApiDB.ApiDB.QuerySpecificationsBySpecificationsId(order.SpecificationsId)
			if err != nil {
				return nil, err
			}
			apiKey = &tables.APIKey{
				ApiKey:       "",
				OrderId:      order.OrderId,
				ApiId:        order.ApiId,
				RequestLimit: spec.Amount,
				UsedNum:      0,
			}
		}
		res[i] = &common.OrderResult{
			Title:        order.Title,
			OrderId:      order.OrderId,
			CreateTime:   order.OrderTime,
			TxHash:       order.TxHash,
			ApiId:        order.ApiId,
			RequestLimit: apiKey.RequestLimit,
			UsedNum:      apiKey.UsedNum,
			Status:       order.OrderStatus,
			ApiKey:       apiKey.ApiKey,
			Price:        order.Price,
			Coin:         order.Coin,
		}
	}
	return map[string]interface{}{
		"total":     total,
		"orderList": res,
	}, nil
}

func (this *SagaOrder) GetQrCodeByOrderId(orderId string) (*common.QrCodeResponse, error) {
	code, err := dao.DefSagaApiDB.OrderDB.QueryQrCodeByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, nil
	}
	return common.BuildQrCodeResult(code.QrCodeId), nil
}

func (this *SagaOrder) GetQrCodeDataById(id string) (*tables.QrCode, error) {
	return dao.DefSagaApiDB.OrderDB.QueryQrCodeByQrCodeId(id)
}
func (this *SagaOrder) GetQrCodeResultById(id string) (string, error) {
	return dao.DefSagaApiDB.OrderDB.QueryQrCodeResultByQrCodeId(id)
}

func (this *SagaOrder) CancelOrder(param *common.OrderIdParam) error {
	return dao.DefSagaApiDB.OrderDB.UpdateOrderStatus(param.OrderId, config.Canceled)
}
func (this *SagaOrder) DeleteOrderByOrderId(param *common.OrderIdParam) error {
	return dao.DefSagaApiDB.OrderDB.DeleteOrderByOrderId(param.OrderId)
}

func (this *SagaOrder) GetTxResult(orderId string) (*common.GetOrderResponse, error) {
	order, err := dao.DefSagaApiDB.OrderDB.QueryOrderByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	res := &common.GetOrderResponse{
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
