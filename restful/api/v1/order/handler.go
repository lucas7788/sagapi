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
	err = core.DefSagaApi.TakeOrder(param)
	if err != nil {
		log.Errorf("[TakeOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func PayOrder(c *gin.Context) {
}
