package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.Use(jwt.JWT())
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/deleteOrder", DeleteOrder)
	orderRouteGroup.POST("/generateTestKey", GenerateTestKey)
	orderRouteGroup.POST("/testAPIKey/:apiKey", TestAPIKey)
	orderRouteGroup.GET("/getQrCodeByOrderId/:orderId", GetQrCodeByOrderId)
	orderRouteGroup.GET("/getQrCodeResultByQrCodeId/:qrCodeId", GetQrCodeResultByQrCodeId)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", GetTxResult)
	orderRouteGroup.GET("/queryOrderByPage/:pageNum/:pageSize", QueryOrderByPage)
	orderRouteGroupNoToken := parent.Group("/onto")
	orderRouteGroupNoToken.POST("/sendTx", SendTx)
	orderRouteGroupNoToken.GET("/getQrCodeDataByQrCodeId/:qrCodeId", GetQrCodeDataByQrCodeId)
}
