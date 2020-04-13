package nasa

import (
	"fmt"
	"github.com/ontio/saga/config"
	"github.com/ontio/saga/core/http"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

func Apod() ([]byte, error) {
	url := fmt.Sprintf(apod, config.NASA_API_KEY)
	return http.Get(url)
}

func Feed(startDate, endDate string) ([]byte, error) {
	url := fmt.Sprintf(feed, startDate, endDate, config.NASA_API_KEY)
	return http.Get(url)
}
