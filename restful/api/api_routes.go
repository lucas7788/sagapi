package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/restful/api/v1"
	"github.com/ontio/sagapi/restful/publish"
)

func RoutesApi(parent *gin.Engine) {
	apiRoutes := parent.Group("/api")
	v1.RoutesV1(apiRoutes)
	publish.RoutesPublish(apiRoutes)
}
