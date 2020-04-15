package order

import (
	"github.com/gin-gonic/gin"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/getQrCodeByOrderId", GetQrCodeByOrderId)
	orderRouteGroup.POST("/getQrCodeDataByQrCodeId", GetQrCodeDataByQrCodeId)
	orderRouteGroup.POST("/sendTx", SendTx)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/deleteOrder", DeleteOrder)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", GetTxResult)
}
