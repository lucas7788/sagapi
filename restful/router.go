package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/middleware/cors"
	"github.com/ontio/sagapi/restful/api"
)

func NewRouter() *gin.Engine {
	root := gin.New()
	root.Use(cors.Cors())
	api.RoutesApi(root)
	return root
}
