package nasa

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core/nasa"
)

func Apod(c *gin.Context) {
	apikey := c.Param("apikey")
	fmt.Printf("apikey: %s\n", apikey)

	res, err := nasa.Apod()
	if err != nil {
		common.WriteResponse(c, nil, common.INTER_ERROR, err)
	} else {
		common.WriteResponse(c, res, common.SUCCESS, nil)
	}
}

func Feed(c *gin.Context) {
	startdate := c.Param("startdate")
	enddate := c.Param("enddate")

	fmt.Printf("startdate: %s, enddate: %s\n", startdate, enddate)
	res, err := nasa.Feed(startdate, enddate)

	if err != nil {
		common.WriteResponse(c, nil, common.INTER_ERROR, err)
	} else {
		common.WriteResponse(c, res, common.SUCCESS, nil)
	}
}
