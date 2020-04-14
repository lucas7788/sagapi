package order

import (
	"fmt"
	"github.com/candybox-sig/log"
	"github.com/gin-gonic/gin"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/restful/api/common"
)

func TakeOrder(c *gin.Context) {
	param := &common2.TakeOrderParam{}
	err := common.ParsePostParam(c.Request.Body, param)
	if err != nil {
		log.Errorf("[TakeOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	fmt.Println(param)
	code, err := core.DefSagaOrder.TakeOrder(param)
	if err != nil {
		log.Errorf("[TakeOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(code))
}

func PayOrder(c *gin.Context) {
}

func SendTx(c *gin.Context) {
	param := &common2.SendTxParam{}
	err := common.ParsePostParam(c.Request.Body, param)
	if err != nil {
		log.Errorf("[SendTx] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	err = core.SendTX(param.SignedTx)
	if err != nil {
		log.Errorf("[SendTx] SendTX failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}
