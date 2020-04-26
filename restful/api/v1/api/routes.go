package api

import "github.com/gin-gonic/gin"

func RoutesApiList(parent *gin.RouterGroup) {
	apiRouteGroup := parent.Group("/apilist")

	apiRouteGroup.POST("/searchApiByKey", SearchApiByKey)

	apiRouteGroup.GET("/getBasicApiInfoByPage/:pageNum/:pageSize", GetBasicApiInfoByPage)
	apiRouteGroup.GET("/getApiDetailByApiId/:apiId", GetApiDetailByApiId)
	apiRouteGroup.GET("/searchApiByCategory/:categoryId", SearchApiByCategoryId)
	apiRouteGroup.GET("/searchApi", SearchApi) //get newest
}
