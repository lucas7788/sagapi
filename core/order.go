package core

import (
	"fmt"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/ontio/sagapi/utils"
	"math/big"
	"strings"
	"sync"
	"time"
)

type SagaOrder struct {
	qrCodeCache *sync.Map //qrCodeId -> QrCode
}

func NewSagaOrder() *SagaOrder {
	return &SagaOrder{
		qrCodeCache: new(sync.Map),
	}
}

func (this *SagaOrder) TakeOrder(param *common.TakeOrderParam) (*common.QrCodeResponse, error) {
	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		return nil, errl
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	info, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(tx, param.ApiId)
	if err != nil {
		errl = err
		return nil, err
	}

	spec, err := dao.DefSagaApiDB.QuerySpecificationsById(tx, param.SpecificationsId)
	if err != nil {
		errl = err
		return nil, err
	}

	price := utils.ToIntByPrecise(spec.Price, sagaconfig.ONG_DECIMALS)
	specifications := new(big.Int).SetUint64(uint64(spec.Amount))
	amount := new(big.Int).Mul(price, specifications)
	amountStr := utils.ToStringByPrecise(amount, sagaconfig.ONG_DECIMALS)
	orderId := common.GenerateUUId(common.UUID_TYPE_RAW)
	order := &tables.Order{
		OrderId:          orderId,
		Title:            info.Title,
		ProductName:      param.ProductName,
		OrderType:        sagaconfig.Api,
		OrderTime:        time.Now().Unix(),
		OrderStatus:      sagaconfig.Processing,
		Amount:           amountStr,
		OntId:            param.OntId,
		UserName:         param.UserName,
		Price:            spec.Price,
		ApiId:            info.ApiId,
		ApiUrl:           info.ApiUrl,
		SpecificationsId: param.SpecificationsId,
		Coin:             sagaconfig.TOKEN_TYPE_ONG,
	}
	err = dao.DefSagaApiDB.InsertOrder(tx, order)
	if err != nil {
		errl = err
		return nil, err
	}
	arr := strings.Split(param.OntId, ":")
	if len(arr) < 3 {
		errl = fmt.Errorf("error ontid: %s", param.OntId)
		return nil, errl
	}
	code, err := common.BuildQrCode(orderId, param.OntId, arr[2], arr[2], sagaconfig.DefSagaConfig.Collect_Money_Address, amountStr)
	if err != nil {
		errl = err
		return nil, err
	}
	//this.qrCodeCache.Store(code.QrCodeId, code)
	err = dao.DefSagaApiDB.InsertQrCode(tx, code)
	if err != nil {
		errl = err
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		errl = err
		return nil, err
	}

	return common.BuildQrCodeResponse(code.QrCodeId), nil
}

func (this *SagaOrder) QueryOrderByPage(pageNum, pageSize int, ontid string) (map[string]interface{}, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 0 {
		pageSize = 0
	}
	total, err := dao.DefSagaApiDB.QueryOrderSum(nil, ontid)
	if err != nil {
		return nil, err
	}
	start := (pageNum - 1) * pageSize
	orders, err := dao.DefSagaApiDB.QueryOrderByPage(nil, start, pageSize, ontid)
	if err != nil {
		return nil, err
	}
	res := make([]*common.OrderResult, len(orders))
	for i, order := range orders {
		apiKey, err := dao.DefSagaApiDB.QueryApiKeyByOrderId(nil, order.OrderId)
		if err != nil && !dao.IsErrNoRows(err) {
			return nil, err
		}
		if apiKey == nil {
			spec, err := dao.DefSagaApiDB.QuerySpecificationsById(nil, order.SpecificationsId)
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
			Amount:       order.Amount,
			CreateTime:   order.OrderTime,
			TxHash:       order.TxHash,
			ApiId:        order.ApiId,
			ApiUrl:       order.ApiUrl,
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

//every time generate new qrcode
func (this *SagaOrder) GetQrCodeByOrderId(ontId, orderId string) (*common.QrCodeResponse, error) {
	tx, err := dao.DefSagaApiDB.DB.Beginx()
	if err != nil {
		return nil, err
	}

	var errl error
	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	arr := strings.Split(ontId, ":")
	if len(arr) != 3 {
		return nil, fmt.Errorf("error ontId: %s", ontId)
	}
	order, err := dao.DefSagaApiDB.QueryOrderByOrderId(tx, orderId)
	if err != nil {
		errl = err
		return nil, err
	}
	code, err := common.BuildQrCode(orderId, ontId, arr[2], arr[2], sagaconfig.DefSagaConfig.Collect_Money_Address, order.Amount)
	err = dao.DefSagaApiDB.InsertQrCode(tx, code)
	if err != nil {
		errl = err
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		errl = err
		return nil, err
	}
	return common.BuildQrCodeResponse(code.QrCodeId), nil
}

func (this *SagaOrder) GetQrCodeDataById(id string) (*tables.QrCode, error) {
	return dao.DefSagaApiDB.QueryQrCodeByQrCodeId(nil, id)
}
func (this *SagaOrder) GetQrCodeResultById(id string) (string, error) {
	return dao.DefSagaApiDB.QueryQrCodeResultByQrCodeId(nil, id)
}

//1. delete qrCodeId
//2. cancel order
func (this *SagaOrder) CancelOrder(orderId string) error {
	tx, err := dao.DefSagaApiDB.DB.Beginx()
	if err != nil {
		return err
	}

	var errl error
	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	status, err := dao.DefSagaApiDB.QueryOrderStatusByOrderId(tx, orderId)
	if err != nil {
		errl = err
		return err
	}
	if status == sagaconfig.Processing {
		//delete qrCodeId
		err = dao.DefSagaApiDB.DeleteQrCodeByOrderId(tx, orderId)
		if err != nil {
			errl = err
			return err
		}
		err = dao.DefSagaApiDB.UpdateOrderStatus(tx, orderId, sagaconfig.Canceled)
		if err != nil {
			errl = err
			return err
		}

		errl = tx.Commit()
		if errl != nil {
			return errl
		}
	}

	errl = fmt.Errorf("only processing order can be canceled")
	return errl
}

//1. delete qrCodeId
//2. cancel order
func (this *SagaOrder) DeleteOrderByOrderId(orderId string) error {
	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		return errl
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	err := dao.DefSagaApiDB.DeleteQrCodeByOrderId(tx, orderId)
	if err != nil {
		errl = err
		return err
	}
	err = dao.DefSagaApiDB.DeleteOrderByOrderId(tx, orderId)
	if err != nil {
		errl = err
		return err
	}
	err = tx.Commit()
	if err != nil {
		errl = err
		return err
	}
	return nil
}

func (this *SagaOrder) GetTxResult(orderId string) (*common.GetOrderResponse, error) {
	order, err := dao.DefSagaApiDB.QueryOrderByOrderId(nil, orderId)
	if err != nil {
		return nil, err
	}
	res := &common.GetOrderResponse{
		UserName: order.UserName,
		OntId:    order.OntId,
	}
	if order.OrderStatus == sagaconfig.Completed {
		res.Result = "1"
	} else if order.OrderStatus == sagaconfig.Processing {
		res.Result = ""
	} else if order.OrderStatus == sagaconfig.Failed {
		res.Result = "0"
	}
	return res, nil
}
