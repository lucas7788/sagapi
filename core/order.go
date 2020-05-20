package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ontio/ontology/common/log"
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

func (this *SagaOrder) TakeWetherForcastApiOrder(param *common.WetherForcastRequest, ontId string) (*common.QrCodeResponse, error) {
	log.Debugf("%v", *param)
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	log.Debugf("%s", string(data))

	toolbox, err := dao.DefSagaApiDB.QueryToolBoxById(nil, uint32(param.ToolBoxId))
	if err != nil {
		log.Debugf("TakeWetherForcastApiOrder. 0 : %s", err)
		return nil, err
	}
	api, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, param.ApiSourceId, tables.API_STATE_BUILTIN)
	if err != nil {
		log.Debugf("TakeWetherForcastApiOrder. 1 : %s", err)
		return nil, err
	}
	if api.ApiType != toolbox.Title || api.ApiKind != tables.API_KIND_DATA_PROCESS {
		log.Debugf("TakeWetherForcastApiOrder. 2 : %s", err)
		return nil, errors.New("wrong api type or kind")
	}

	// to do check the conrespond.
	alg, err := dao.DefSagaApiDB.QueryAlgorithmById(nil, param.AlgorithmId)
	if err != nil {
		log.Debugf("TakeWethe2ForcastApiOrder. 3 : %s", err)
		return nil, err
	}
	env, err := dao.DefSagaApiDB.QueryEnvById(nil, param.EnvId)
	if err != nil {
		log.Debugf("TakeWethe2ForcastApiOrder. 4 : %s", err)
		return nil, err
	}

	tx, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("TakeWethe2ForcastApiOrder. 4 : %s", err)
		return nil, errl
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	apiprice := utils.ToIntByPrecise(api.Price, sagaconfig.ONG_DECIMALS)
	algprice := utils.ToIntByPrecise(alg.Price, sagaconfig.ONG_DECIMALS)
	envprice := utils.ToIntByPrecise(env.Price, sagaconfig.ONG_DECIMALS)
	t0 := big.NewInt(0)
	t1 := t0.Add(apiprice, algprice)
	t2 := big.NewInt(0)
	allAmount := t2.Add(t1, envprice)
	amountStr := utils.ToStringByPrecise(allAmount, sagaconfig.ONG_DECIMALS)

	orderId := common.GenerateUUId(common.UUID_TYPE_ORDER_ID)
	order := &tables.Order{
		OrderId:     orderId,
		ApiId:       api.ApiId,
		Title:       toolbox.Title,
		OrderType:   sagaconfig.ApiProcess,
		OrderTime:   time.Now().Unix(),
		OrderStatus: sagaconfig.Processing,
		OntId:       ontId,
		Price:       amountStr,
		OrderKind:   tables.ORDER_KIND_DATA_PROCESS_WETHER,
		Coin:        api.Coin,
		Request:     string(data), // leave it empty fill later.
	}
	log.Debugf("Request: %s", order.Request)

	err = dao.DefSagaApiDB.InsertOrder(tx, order)
	if err != nil {
		errl = err
		return nil, err
	}
	arr := strings.Split(ontId, ":")
	if len(arr) < 3 {
		log.Debugf("TakeWethe2ForcastApiOrder. 5 : %s", err)
		errl = fmt.Errorf("error ontid: %s", ontId)
		return nil, errl
	}
	payer := arr[2]

	code, err := common.BuildWetherForcastQrCode(sagaconfig.DefSagaConfig.NetType, orderId, ontId, api.ResourceId, alg.ResourceId, api.TokenHash, env.TokenHash, payer, env.OwnerAddress, amountStr)
	if err != nil {
		errl = err
		log.Debugf("TakeWethe2ForcastApiOrder. 6 : %s", err)
		return nil, err
	}

	err = dao.DefSagaApiDB.InsertQrCode(tx, code)
	if err != nil {
		errl = err
		log.Debugf("TakeWethe2ForcastApiOrder. 7 : %s", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		errl = err
		return nil, err
	}

	return common.BuildQrCodeResponse(code.QrCodeId), nil
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

	info, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(tx, param.ApiId, tables.API_STATE_BUILTIN)
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
	orderId := common.GenerateUUId(common.UUID_TYPE_ORDER_ID)
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
		OrderKind:        tables.ORDER_KIND_API,
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
	total, err := dao.DefSagaApiDB.QueryOrderSum(nil, ontid, sagaconfig.Api)
	if err != nil {
		return nil, err
	}

	log.Debugf("QueryOrderByPage.Y 0 %d", total)
	start := (pageNum - 1) * pageSize
	orders, err := dao.DefSagaApiDB.QueryOrderByPage(nil, start, pageSize, ontid)
	if err != nil {
		log.Debugf("QueryOrderByPage.N 0 %d", err)
		return nil, err
	}
	res := make([]*common.OrderResult, 0)
	for _, order := range orders {
		if order.OrderType == sagaconfig.Api {
			// todo handle WetherForcast api and data process.
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
			res = append(res, &common.OrderResult{
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
			})
		}
	}
	return map[string]interface{}{
		"total":     total,
		"orderList": res,
	}, nil
}

func (this *SagaOrder) QueryDataProcessOrderByPage(pageNum, pageSize int, ontid string) (map[string]interface{}, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 0 {
		pageSize = 0
	}
	total, err := dao.DefSagaApiDB.QueryOrderSum(nil, ontid, sagaconfig.ApiProcess)
	if err != nil {
		return nil, err
	}

	log.Debugf("QueryOrderByPage.Y 0 %d", total)
	start := (pageNum - 1) * pageSize
	orders, err := dao.DefSagaApiDB.QueryOrderByPage(nil, start, pageSize, ontid)
	if err != nil {
		log.Debugf("QueryOrderByPage.N 0 %d", err)
		return nil, err
	}
	res := make([]*common.DataProcessOrderResult, 0)
	for _, order := range orders {
		if order.OrderType == sagaconfig.ApiProcess {
			// todo handle WetherForcast api and data process.
			res = append(res, &common.DataProcessOrderResult{
				Title:     order.Title,
				OrderId:   order.OrderId,
				OrderTime: order.OrderTime,
				TxHash:    order.TxHash,
				ApiId:     order.ApiId,
				Status:    order.OrderStatus,
				Price:     order.Price,
				Coin:      order.Coin,
				OrderKind: order.OrderKind,
				Request:   order.Request,
				Result:    order.Result,
			})
		}
	}
	return map[string]interface{}{
		"total":     total,
		"orderList": res,
	}, nil
}

func (this *SagaOrder) GetOrderDetailById(orderId string) (*common.OrderDetailResponse, error) {
	order, err := dao.DefSagaApiDB.QueryOrderByOrderId(nil, orderId)
	if err != nil {
		return nil, err
	}

	switch order.OrderKind {
	case tables.ORDER_KIND_DATA_PROCESS_WETHER:
		paramWether := &common.WetherForcastRequest{}
		err = json.Unmarshal([]byte(order.Request), paramWether)
		if err != nil {
			log.Debugf("GetOrderDetailById.N.0: %s", err)
			return nil, err
		}
		toolbox, err := dao.DefSagaApiDB.QueryToolBoxById(nil, uint32(paramWether.ToolBoxId))
		if err != nil {
			log.Debugf("TakeWetherForcastApiOrder. 0 : %s", err)
			return nil, err
		}
		env, err := dao.DefSagaApiDB.QueryEnvById(nil, paramWether.EnvId)
		if err != nil {
			log.Debugf("GetOrderDetailById.N.1: %s", err)
			return nil, err
		}
		alg, err := dao.DefSagaApiDB.QueryAlgorithmById(nil, paramWether.AlgorithmId)
		if err != nil {
			log.Debugf("GetOrderDetailById.N.2: %s", err)
			return nil, err
		}
		api, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, order.ApiId, tables.API_STATE_BUILTIN)
		if err != nil {
			log.Debugf("TakeWetherForcastApiOrder. 1 : %s", err)
			return nil, err
		}

		detail := common.WetherOrderDetail{
			TargetDate: paramWether.TargetDate,
			Location:   &paramWether.Location,
			ToolBox:    toolbox,
			ApiSource:  api,
			Algorithm:  alg,
			Env:        env,
			Result:     order.Result,
		}
		return &common.OrderDetailResponse{
			Res: detail,
		}, nil
	default:
		return nil, fmt.Errorf("wrong order kind %d", order.OrderKind)

	}
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
	var code *tables.QrCode
	if order.OrderType == sagaconfig.Api {
		code, err = common.BuildQrCode(orderId, ontId, arr[2], arr[2], sagaconfig.DefSagaConfig.Collect_Money_Address, order.Amount)
	} else if order.OrderType == sagaconfig.ApiProcess {
		switch order.OrderKind {
		case tables.ORDER_KIND_DATA_PROCESS_WETHER:
			api, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(nil, order.ApiId, tables.API_STATE_BUILTIN)
			if err != nil {
				log.Debugf("GetQrCodeByOrderId.N.0: %s", err)
				return nil, err
			}

			paramWether := &common.WetherForcastRequest{}
			err = json.Unmarshal([]byte(order.Request), paramWether)
			if err != nil {
				errl = err
				log.Debugf("GetQrCodeByOrderId.N.1: %s", err)
				return nil, err
			}
			env, err := dao.DefSagaApiDB.QueryEnvById(nil, paramWether.EnvId)
			if err != nil {
				errl = err
				log.Debugf("GetQrCodeByOrderId.N.2: %s", err)
				return nil, err
			}

			alg, err := dao.DefSagaApiDB.QueryAlgorithmById(nil, paramWether.AlgorithmId)
			if err != nil {
				errl = err
				log.Debugf("GetQrCodeByOrderId.N.3: %s", err)
				return nil, err
			}

			code, err = common.BuildWetherForcastQrCode(sagaconfig.DefSagaConfig.NetType, orderId, ontId, api.ResourceId, alg.ResourceId, api.TokenHash, env.TokenHash, arr[2], env.OwnerAddress, order.Price)
		default:
			return nil, errors.New("wrong order type title.")
		}

	}

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
	} else {
		errl = fmt.Errorf("only processing order can be canceled")
	}

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
