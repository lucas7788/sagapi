package order

import (
	"github.com/gin-gonic/gin"
)

func RoutesOrder(parent *gin.RouterGroup) {
	nasaRouteGroup := parent.Group("/order")
	nasaRouteGroup.POST("/takeOrder", TakeOrder)
	nasaRouteGroup.POST("/payOrder", PayOrder)
}
