package order

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
	"github.com/ontio/sagapi/sagaconfig"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TakeOrder(c *gin.Context) {
	param := &common2.TakeOrderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[TakeOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(sagaconfig.Key_OntId)
	if !ok {
		log.Errorf("[TakeOrder] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	param.OntId = ontid.(string)
	res, err := core.DefSagaApi.SagaOrder.TakeOrder(param)
	if err != nil {
		log.Errorf("[TakeOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func QueryOrderByPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[QueryOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(sagaconfig.Key_OntId)
	if !ok || ontId == "" {
		log.Errorf("[QueryOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	log.Infof("[QueryOrderByPage] ontid:%s", ontId)
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	paseSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	orders, err := core.DefSagaApi.SagaOrder.QueryOrderByPage(pageNum, paseSize, ontId.(string))
	if err != nil {
		log.Errorf("[QueryOrderByPage] QueryOrderByPage failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(orders))
}

func GenerateTestKey(c *gin.Context) {
	params := &common2.GenerateTestKeyParam{}
	err := common.ParsePostParam(c, params)
	if err != nil {
		log.Errorf("[GenerateTestKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(sagaconfig.Key_OntId)
	if !ok {
		log.Errorf("[GenerateTestKey] ontId is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	testKey, err := core.DefSagaApi.GenerateApiTestKey(params.ApiId, ontId.(string))
	if err != nil || testKey == nil {
		log.Errorf("[GenerateTestKey] GenerateApiTestKey failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(testKey))
}

func TestAPIKey(c *gin.Context) {
	var params []tables.RequestParam
	err := common.ParsePostParam(c, &params)
	if err != nil {
		log.Errorf("[GenerateTestKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	data, err := core.DefSagaApi.TestApiKey(params)
	if err != nil {
		log.Errorf("[TestAPIKey] TestApiKey failed: %s", err.Error())
		res := make(map[string]string)
		res["errorDesc"] = err.Error()
		bs, _ := json.Marshal(res)
		common.WriteResponse(c, common.ResponseSuccess(string(bs)))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(string(data)))
}

func GetQrCodeByOrderId(c *gin.Context) {
	paramArr, err := common.ParseGetParamByParamName(c, "orderId")
	if err != nil {
		log.Errorf("[GetQrCodeByOrderId] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(sagaconfig.Key_OntId)
	if !ok || ontId == nil {
		log.Errorf("[GetQrCodeByOrderId] ontid is nil")
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	res, err := core.DefSagaApi.SagaOrder.GetQrCodeByOrderId(ontId.(string), paramArr[0])
	if err != nil {
		log.Errorf("[TakeOrGetQrCodeByOrderIdder] GetQrCodeByOrderId failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

//onto invoke
func GetQrCodeDataByQrCodeId(c *gin.Context) {
	paramArr, err := common.ParseGetParamByParamName(c, "qrCodeId")
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	code, err := core.DefSagaApi.SagaOrder.GetQrCodeDataById(paramArr[0])
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	c.JSON(http.StatusOK, code)
}

func GetQrCodeResultByQrCodeId(c *gin.Context) {
	paramArr, err := common.ParseGetParamByParamName(c, "qrCodeId")
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	code, err := core.DefSagaApi.SagaOrder.GetQrCodeResultById(paramArr[0])
	if err != nil {
		log.Errorf("[SendTx] GetQrCodeResultById failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(code))
}

func CancelOrder(c *gin.Context) {
	param := &common2.OrderIdParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[CancelOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if param == nil || param.OrderId == "" {
		log.Errorf("[CancelOrder] param is nil failed")
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("param is nil")))
		return
	}
	err = core.DefSagaApi.SagaOrder.CancelOrder(param.OrderId)
	if err != nil {
		log.Errorf("[CancelOrder] CancelOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	log.Infof("[CancelOrder] orderId:%s", param.OrderId)
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func DeleteOrder(c *gin.Context) {
	param := &common2.OrderIdParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[CancelOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	err = core.DefSagaApi.SagaOrder.DeleteOrderByOrderId(param.OrderId)
	if err != nil {
		log.Errorf("[CancelOrder] CancelOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func SendTx(c *gin.Context) {
	param := &common2.SendTxParam{}
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	err = json.Unmarshal(paramsBs, param)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	err = core.SendTX(param)
	if err != nil {
		log.Errorf("[SendTx] SendTX failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.INTER_ERROR, err))
		return
	}
	writeOntoResponse(c, common.ResponseSuccessOnto())
}

func writeOntoResponse(c *gin.Context, param map[string]interface{}) {
	c.JSON(http.StatusOK, param)
}

func GetTxResult(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := core.DefSagaApi.SagaOrder.GetTxResult(orderId)
	if err != nil {
		log.Errorf("[GetTxResult] QueryOrderByOrderId failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}
