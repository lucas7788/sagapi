package nasa

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/restful/api/common"
)

func Apod(c *gin.Context) {
	apikey := c.Param("apikey")
	if apikey == "" {
		log.Errorf("[nasa_handler] apod error: %s")
		common.WriteResponse(c, common.ResponseFailed(common.API_KEY_IS_NIL, fmt.Errorf("api key is nil")))
		return
	}
	res, err := core.DefSagaApi.Nasa.Apod(apikey)
	if err != nil {
		log.Errorf("[nasa_handler] apod error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(string(res)))
}

func Feed(c *gin.Context) {
	startdate := c.Param("startdate")
	enddate := c.Param("enddate")
	apikey := c.Param("apikey")

	fmt.Printf("startdate: %s, enddate: %s\n", startdate, enddate)
	res, err := core.DefSagaApi.Nasa.Feed(startdate, enddate, apikey)
	if err != nil {
		log.Errorf("[nasa_handler] apod error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(string(res)))
}
