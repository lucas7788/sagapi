package publish

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/jwt"
)

func RoutesPublish(parent *gin.Engine) {
	publishG := parent.Group("/publish")
	publishG.POST("/api", PublishAPIHandle)

	publishGAmin := publishG.Group("/admin")
	publishGAmin.Use(jwt.JWTAdmin())
	publishGAmin.POST("/publish", VerifyAPIHandle)
	publishGAmin.POST("/getallpublishapi/:pageNum/:pageSize", GetALLPublishPage)
}
