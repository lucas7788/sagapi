package nasa

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/core/nasa"
)

func Apod(c *gin.Context) {
	apikey := c.Param("apikey")
	fmt.Printf("apikey: %s\n", apikey)

	res, err := nasa.Apod()
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.String(http.StatusOK, string(res))
	}
}

func Feed(c *gin.Context) {
	startdate := c.Param("startdate")
	enddate := c.Param("enddate")

	fmt.Printf("startdate: %s, enddate: %s\n", startdate, enddate)
	res, err := nasa.Feed(startdate, enddate)
	if err != nil {
		c.String(http.StatusNotFound, "")
	} else {
		c.String(http.StatusOK, string(res))
	}
}
