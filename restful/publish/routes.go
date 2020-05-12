package publish

import (
	"github.com/gin-gonic/gin"
)

func RoutesPublish(parent *gin.Engine) {
	publishG := parent.Group("/publish")
	publishG.POST("/api", PublishAPIHandle)
	publishG.POST("/publish", VerifyAPIHandle)
	publishG.POST("/getallpublishapi/:pageNum/:pageSize", GetALLPublishPage)
}
