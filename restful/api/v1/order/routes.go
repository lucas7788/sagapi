package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.Use(jwt.JWT())
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/sendTx", SendTx)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/deleteOrder", DeleteOrder)
	orderRouteGroup.POST("/generateTestKey", GenerateTestKey)
	orderRouteGroup.GET("/getQrCodeByOrderId/:orderId", GetQrCodeByOrderId)
	orderRouteGroup.GET("/getQrCodeDataByQrCodeId/:qrCodeId", GetQrCodeDataByQrCodeId)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", GetTxResult)
}
