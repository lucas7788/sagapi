package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/restful/api/v1/nasa"
)

func RoutesV1(parent *gin.RouterGroup) {
	v1Route := parent.Group("/v1")
	nasa.RoutesNasa(v1Route)
	parent.POST("/takeOrder", TakeOrder)
	parent.POST("/payOrder", PayOrder)
}
