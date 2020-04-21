package api

import "github.com/gin-gonic/gin"

func RoutesApiList(parent *gin.RouterGroup) {
	apiRouteGroup := parent.Group("/apilist")
	apiRouteGroup.GET("/getBasicApiInfoByPage/:pageNum/:pageSize", GetBasicApiInfoByPage)
	apiRouteGroup.GET("/getApiDetailByApiId/:apiId", GetApiDetailByApiId)
	apiRouteGroup.GET("/searchApiByKey/:key", SearchApiByKey)
	apiRouteGroup.GET("/searchApiByCategory/:categoryId", SearchApiByCategoryId)
	apiRouteGroup.GET("/searchApi", SearchApi) //get newest
}
