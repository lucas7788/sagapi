package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/restful/api"
)

func NewRouter() *gin.Engine {
	root := gin.New()
	api.RoutesApi(root)
	return root
}
