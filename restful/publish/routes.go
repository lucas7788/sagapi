package publish

import (
	"github.com/gin-gonic/gin"
)

func RoutesPublish(parent *gin.Engine) {
	publishG := parent.Group("/publish")
	publishG.POST("/api", PublishAPIHandle)
	publishG.POST("/verify", VerifyAPIHandle)
	publishG.POST("/getallpublishapi", GetALLPublishPage)
}
