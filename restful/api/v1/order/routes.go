package order

import (
	"github.com/gin-gonic/gin"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/getQrCode", GetQrCodeById)
	orderRouteGroup.POST("/sendTx", SendTx)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/deleteOrder", DeleteOrder)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", QueryOrderStatus)
}
