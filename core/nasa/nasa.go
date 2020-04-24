package nasa

import (
	"errors"
	"fmt"

	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"strings"
	"sync"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

type ApiKeyInvokeFre struct {
	ApiKey    *tables.APIKey
	InvokeFre int
}

type Nasa struct {
	apiKeyInvokeFreCache *sync.Map //apikey -> ApiKeyInvokeFre
}

func NewNasa() *Nasa {
	return &Nasa{
		apiKeyInvokeFreCache: new(sync.Map),
	}
}

func (this *Nasa) beforeCheckApiKey(apiKey string) (*ApiKeyInvokeFre, error) {
	key, err := this.getApiKeyInvokeFre(apiKey)
	if err != nil {
		return nil, err
	}
	if key.ApiKey.UsedNum >= key.ApiKey.RequestLimit {
		return nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.ApiKey.UsedNum, key.ApiKey.RequestLimit)
	}
	return key, nil
}

func (this *Nasa) ApodParams(params []tables.RequestParam) (string, error) {
	if len(params) == 1 && params[0].ParamName == "apiKey" {
		return this.Apod(params[0].ValueDesc)
	}
	return "", errors.New("Apod params error")
}

func (this *Nasa) FeedParams(params []tables.RequestParam) (string, error) {
	if len(params) == 3 && params[0].ParamName == "startDate" && params[1].ParamName == "endDate" && params[2].ParamName == "apiKey" {
		return this.Feed(params[0].ValueDesc, params[1].ValueDesc, params[2].ValueDesc)
	}
	return "", errors.New("Apod params error")
}

func (this *Nasa) Apod(apiKey string) (string, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf(apod, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		return "", err
	}
	key.ApiKey.UsedNum += 1

	//TODO
	err = this.updateApiKeyInvokeFre(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) (string, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf(feed, startDate, endDate, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		return "", err
	}
	//TODO
	err = this.updateApiKeyInvokeFre(key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (this *Nasa) getApiKeyInvokeFre(apiKey string) (*ApiKeyInvokeFre, error) {
	keyIn, ok := this.apiKeyInvokeFreCache.Load(apiKey)
	if !ok || keyIn == nil {
		var err error
		var key *tables.APIKey
		if isTestKey(apiKey) {
			key, err = dao.DefSagaApiDB.ApiDB.QueryApiTestKeyByApiTestKey(apiKey)
		} else {
			key, err = dao.DefSagaApiDB.ApiDB.QueryApiKeyByApiKey(apiKey)
		}
		if err != nil {
			return nil, err
		}
		invokeFre, err := dao.DefSagaApiDB.ApiDB.QueryInvokeFreByApiId(key.ApiId)
		if err != nil {
			return nil, err
		}
		return &ApiKeyInvokeFre{
			ApiKey:    key,
			InvokeFre: invokeFre,
		}, nil
	} else {
		return keyIn.(*ApiKeyInvokeFre), nil
	}
}

func isTestKey(apiKey string) bool {
	return strings.Contains(apiKey, "test")
}

func (this *Nasa) updateApiKeyInvokeFre(key *ApiKeyInvokeFre) error {
	this.apiKeyInvokeFreCache.Store(key.ApiKey, key)
	if strings.Contains(key.ApiKey.ApiKey, "test") {
		err := dao.DefSagaApiDB.ApiDB.UpdateApiTestKeyUsedNum(key.ApiKey.ApiKey, key.ApiKey.UsedNum)
		if err != nil {
			return err
		}
	} else {
		err := dao.DefSagaApiDB.ApiDB.UpdateApiKeyUsedNum(key.ApiKey.ApiKey, key.ApiKey.UsedNum)
		if err != nil {
			return err
		}
	}
	return dao.DefSagaApiDB.ApiDB.UpdateInvokeFrequencyByApiId(key.InvokeFre, key.ApiKey.ApiId)
}
