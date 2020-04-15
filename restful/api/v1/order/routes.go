package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.Use(jwt.JWT())
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/getQrCodeByOrderId", GetQrCodeByOrderId)
	orderRouteGroup.POST("/getQrCodeDataByQrCodeId", GetQrCodeDataByQrCodeId)
	orderRouteGroup.POST("/sendTx", SendTx)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/deleteOrder", DeleteOrder)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", GetTxResult)
}
