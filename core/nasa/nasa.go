package nasa

import (
	"errors"
	"fmt"
	"github.com/ontio/sagapi/core/freq"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"sync/atomic"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

type Nasa struct {
	Cache *freq.DBCache
}

func NewNasa(cache *freq.DBCache) *Nasa {
	res := &Nasa{
		Cache: cache,
	}

	return res
}

func (this *Nasa) Apod(apiKey string) ([]byte, error) {
	key, apiCounterP, err := this.Cache.BeforeCheckApiKey(apiKey, sagaconfig.APOD)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(apod, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		atomic.AddUint64(&key.UsedNum, ^uint64(0))
		atomic.AddUint64(apiCounterP, ^uint64(0))
		return nil, err
	}

	this.Cache.UpdateFreq <- apiKey

	return res, nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) ([]byte, error) {
	key, apiCounterP, err := this.Cache.BeforeCheckApiKey(apiKey, sagaconfig.FEED)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(feed, startDate, endDate, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		atomic.AddUint64(&key.UsedNum, ^uint64(0))
		atomic.AddUint64(apiCounterP, ^uint64(0))
		return nil, err
	}

	this.Cache.UpdateFreq <- apiKey
	return res, nil
}

func (this *Nasa) ApodParams(apiKey string) ([]byte, error) {
	return this.Apod(apiKey)
}

func (this *Nasa) FeedParams(params []*tables.RequestParam, apiKey string) ([]byte, error) {
	if len(params) == 2 && params[0].ParamName == "startDate" && params[1].ParamName == "endDate" {
		return this.Feed(params[0].ValueDesc, params[1].ValueDesc, apiKey)
	}
	return nil, errors.New("Apod params error")
}
