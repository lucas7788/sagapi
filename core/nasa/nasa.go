package nasa

import (
	"errors"
	"fmt"

	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"strings"
	"sync"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

type Nasa struct {
	apiKeyCache *sync.Map //apikey -> ApiKey
	invokeFre   *sync.Map //apiId -> invokeFre
}

func NewNasa() *Nasa {
	return &Nasa{
		apiKeyCache: new(sync.Map),
		invokeFre:   new(sync.Map),
	}
}

func (this *Nasa) beforeCheckApiKey(apiKey string) (*tables.APIKey, error) {
	key, err := this.getApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	if key.UsedNum >= key.RequestLimit {
		return nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.RequestLimit)
	}
	return key, nil
}

func (this *Nasa) ApodParams(params []tables.RequestParam) ([]byte, error) {
	if len(params) == 1 && params[0].ParamName == "apiKey" {
		return this.Apod(params[0].ValueDesc)
	}
	return nil, errors.New("Apod params error")
}

func (this *Nasa) FeedParams(params []tables.RequestParam) ([]byte, error) {
	if len(params) != 3 && params[0].ParamName == "startDate" && params[1].ParamName == "endDate" && params[2].ParamName == "apiKey" {
		return this.Feed(params[0].ValueDesc, params[1].ValueDesc, params[2].ValueDesc)
	}
	return nil, errors.New("Apod params error")
}

func (this *Nasa) Apod(apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(apod, config.DefConfig.NASAAPIKey)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	key.UsedNum += 1

	//TODO
	err = this.updateApiKey(key)
	if err != nil {
		return nil, err
	}
	err = this.updateInvokeFreByApiId(key.ApiId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(feed, startDate, endDate, config.DefConfig.NASAAPIKey)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//TODO
	err = this.updateApiKey(key)
	if err != nil {
		return nil, err
	}
	err = this.updateInvokeFreByApiId(key.ApiId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *Nasa) getInvokeFreByApiId(apiId int) (int, error) {
	val, ok := this.invokeFre.Load(apiId)
	if ok {
		return val.(int), nil
	}
	invokeFre, err := dao.DefSagaApiDB.ApiDB.QueryInvokeFreByApiId(apiId)
	if err != nil {
		return 0, err
	}
	return invokeFre, nil
}

func (this *Nasa) updateInvokeFreByApiId(apiId int) error {
	invokeFre, err := this.getInvokeFreByApiId(apiId)
	if err != nil {
		return err
	}
	invokeFre += 1
	return dao.DefSagaApiDB.ApiDB.UpdateInvokeFrequencyByApiId(invokeFre, apiId)
}

func (this *Nasa) getApiKey(apiKey string) (*tables.APIKey, error) {
	keyIn, ok := this.apiKeyCache.Load(apiKey)
	var key *tables.APIKey
	if !ok || keyIn == nil {
		var err error
		if strings.Contains(apiKey, "test") {
			key, err = dao.DefSagaApiDB.ApiDB.QueryApiTestKeyByApiTestKey(apiKey)
		} else {
			key, err = dao.DefSagaApiDB.ApiDB.QueryApiKeyByApiKey(apiKey)
		}
		if err != nil {
			return nil, err
		}
		return key, nil
	} else {
		key = keyIn.(*tables.APIKey)
	}
	return key, nil
}

func (this *Nasa) updateApiKey(key *tables.APIKey) error {
	this.apiKeyCache.Store(key.ApiKey, key)
	if strings.Contains(key.ApiKey, "test") {
		err := dao.DefSagaApiDB.ApiDB.UpdateApiTestKeyUsedNum(key.ApiKey, key.UsedNum)
		if err != nil {
			return err
		}
	} else {
		err := dao.DefSagaApiDB.ApiDB.UpdateApiKeyUsedNum(key.ApiKey, key.UsedNum)
		if err != nil {
			return err
		}
	}
	return nil
}
