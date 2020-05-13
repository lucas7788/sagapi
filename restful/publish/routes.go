package publish

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesPublish(parent *gin.RouterGroup) {
	publishG := parent.Group("/publish")
	publishG.Use(jwt.JWT())
	publishG.POST("/api", PublishAPIHandle)
	publishG.GET("/getpublishapi/:pageNum/:pageSize", GetPulishApi)

	publishGAmin := publishG.Group("/admin")
	publishGAmin.Use(jwt.JWTAdmin())
	publishGAmin.GET("/getallpublishapi/:pageNum/:pageSize", GetALLPublishPage)
	publishGAmin.GET("/getapidetailinfo/:apiId/:apiState", GetApiDetailByApiIdApiState)

	publishGAmin.POST("/admintest/:apiId", AdminTestAPIKey)
	publishGAmin.POST("/publish/:apiId/:sagaUrlKey", VerifyAPIHandle)
}
