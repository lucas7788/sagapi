package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/restful/api/v1/nasa"
	"github.com/ontio/sagapi/restful/api/v1/order"
)

func RoutesV1(parent *gin.RouterGroup) {
	v1Route := parent.Group("/v1")
	nasa.RoutesNasa(v1Route)
	order.RoutesOrder(v1Route)
}
