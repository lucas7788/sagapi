package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/cors"
	"github.com/ontio/sagapi/restful/api/v1/api"
	"github.com/ontio/sagapi/restful/api/v1/nasa"
	"github.com/ontio/sagapi/restful/api/v1/order"
)

func RoutesV1(parent *gin.RouterGroup) {
	v1Route := parent.Group("/v1")
	v1Route.Use(cors.Cors())
	nasa.RoutesNasa(v1Route)
	order.RoutesOrder(v1Route)
	api.RoutesApiList(v1Route)

	v1Route.POST("/data_source/:sagaUrlKey/:apikey", HandleDataSourceReq)
}
