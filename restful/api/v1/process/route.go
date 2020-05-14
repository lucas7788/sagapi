package process

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesDataProcess(parent *gin.RouterGroup) {
	publishG := parent.Group("/dataprocess")
	publishG.Use(jwt.JWT())
	publishG.GET("/GetWetherForcastInfo/:preditype", GetWetherForcastInfo)
	publishG.GET("/GetLocation/:country", GetLocation)
}
