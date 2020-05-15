package process

import (
	"github.com/gin-gonic/gin"
)

func RoutesDataProcess(parent *gin.RouterGroup) {
	publishG := parent.Group("/dataprocess")
	//publishG.Use(jwt.JWT())
	publishG.GET("/GetWetherForcastInfo/:preditype", GetWetherForcastInfo)
	publishG.GET("/GetLocation/:country", GetLocation)
}
