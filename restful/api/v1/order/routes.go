package order

import (
	"github.com/gin-gonic/gin"
)

func RoutesOrder(parent *gin.RouterGroup) {
	nasaRouteGroup := parent.Group("/nasa")
	nasaRouteGroup.POST("/order/takeOrder", TakeOrder)
	nasaRouteGroup.GET("/order/payOrder", PayOrder)
}
