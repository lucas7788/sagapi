package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/restful/api/common"
)

func TakeOrder(c *gin.Context) {
	param := &common2.TakeOrderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[TakeOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	res, err := core.DefSagaOrder.TakeOrder(param)
	if err != nil {
		log.Errorf("[TakeOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetQrCodeByOrderId(c *gin.Context) {
	param := &common2.OrderIdParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[GetQrCodeByOrderId] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	res, err := core.DefSagaOrder.GetQrCodeByOrderId(param.OrderId)
	if err != nil {
		log.Errorf("[TakeOrGetQrCodeByOrderIdder] GetQrCodeByOrderId failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetQrCodeDataByQrCodeId(c *gin.Context) {
	param := &common2.GetQrCodeParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	code, err := core.DefSagaOrder.GetPayQrCodeById(param.Id)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
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
	err = core.DefSagaOrder.CancelOrder(param)
	if err != nil {
		log.Errorf("[CancelOrder] CancelOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
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
	err = core.DefSagaOrder.DeleteOrderByOrderId(param)
	if err != nil {
		log.Errorf("[CancelOrder] CancelOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func SendTx(c *gin.Context) {
	param := &common2.SendTxParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	err = core.SendTX(param)
	if err != nil {
		log.Errorf("[SendTx] SendTX failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess("SUCCESS"))
}

func GetTxResult(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := core.DefSagaOrder.GetTxResult(orderId)
	if err != nil {
		log.Errorf("[GetTxResult] QueryOrderByOrderId failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}
