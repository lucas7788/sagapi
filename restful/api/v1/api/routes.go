package api

import "github.com/gin-gonic/gin"

func RoutesApiList(parent *gin.RouterGroup) {
	apiRouteGroup := parent.Group("/apilist")
	apiRouteGroup.GET("/getBasicApiInfoByPage/:pageNum/:pageSize", GetBasicApiInfoByPage)
	apiRouteGroup.GET("/getApiDetailByApiId/:apiId", GetApiDetailByApiId)
}
