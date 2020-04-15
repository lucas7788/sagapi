package nasa

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesNasa(parent *gin.RouterGroup) {
	nasaRouteGroup := parent.Group("/nasa")
	nasaRouteGroup.Use(jwt.JWT())
	nasaRouteGroup.GET("/apod/:apikey", Apod)
	nasaRouteGroup.GET("/feed/:startdate/:enddate/:apikey", Feed)
}
