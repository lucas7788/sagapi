package nasa

import (
	"errors"
	"fmt"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"sync"
	"sync/atomic"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

type Nasa struct {
	apiKeyCache *sync.Map //apikey -> ApiKeyInvokeFre
	freqLock    *sync.Mutex
	updateFreq  chan string
	invokeFreq  int32
}

func NewNasa() *Nasa {
	res := &Nasa{
		apiKeyCache: new(sync.Map),
		freqLock:    new(sync.Mutex),
		updateFreq:  make(chan string, 20),
	}

	go res.UpdateFreqDataBase()
	return res
}

func (this *Nasa) UpdateFreqDataBase() {
	for {
		select {
		case apiKey := <-this.updateFreq:
			keyIn, ok := this.apiKeyCache.Load(apiKey)
			if !ok {
				continue
			}

			key := keyIn.(*tables.APIKey)
			this.updateApiKeyInvokeFre(key)
			dao.DefSagaApiDB.ApiDB.UpdateApiKeyInvokeFre(key.ApiKey, key.ApiId, key.UsedNum, this.invokeFreq)
		}
	}
}

func (this *Nasa) beforeCheckApiKey(apiKey string, apiId int) (*tables.APIKey, error) {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyInvokeFre(apiKey)
	if err != nil {
		return nil, err
	}

	if key.UsedNum >= key.RequestLimit {
		return nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.RequestLimit)
	}
	if key.ApiId != apiId {
		return nil, fmt.Errorf("this apikey: %s can not invoke this api", apiKey)
	}
	key.UsedNum += 1
	this.invokeFreq += 1
	return key, nil
}

func (this *Nasa) Apod(apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey, sagaconfig.APOD)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(apod, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		atomic.AddInt32(&key.UsedNum, -1)
		atomic.AddInt32(&this.invokeFreq, -1)
		return nil, err
	}

	this.updateFreq <- apiKey

	return res, nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey, sagaconfig.FEED)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(feed, startDate, endDate, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, err := http.DefClient.Get(url)
	if err != nil {
		atomic.AddInt32(&key.UsedNum, -1)
		atomic.AddInt32(&this.invokeFreq, -1)
		return nil, err
	}

	this.updateFreq <- apiKey
	return res, nil
}

func (this *Nasa) ApodParams(params []tables.RequestParam) ([]byte, error) {
	if len(params) == 1 && params[0].ParamName == "apiKey" {
		return this.Apod(params[0].ValueDesc)
	}
	return nil, errors.New("Apod params error")
}

func (this *Nasa) FeedParams(params []tables.RequestParam) ([]byte, error) {
	if len(params) == 3 && params[0].ParamName == "startDate" && params[1].ParamName == "endDate" && params[2].ParamName == "apiKey" {
		return this.Feed(params[0].ValueDesc, params[1].ValueDesc, params[2].ValueDesc)
	}
	return nil, errors.New("Apod params error")
}

func (this *Nasa) getApiKeyInvokeFre(apiKey string) (*tables.APIKey, error) {
	keyIn, ok := this.apiKeyCache.Load(apiKey)
	if !ok || keyIn == nil {
		key, err := dao.DefSagaApiDB.ApiDB.QueryApiKeyByApiKey(apiKey)
		if err != nil {
			return nil, err
		}
		this.apiKeyCache.Store(apiKey, key)
		freq, err := dao.DefSagaApiDB.ApiDB.QueryInvokeFreByApiId(key.ApiId)
		if err != nil {
			return nil, err
		}
		this.invokeFreq = freq
		return key, nil
	} else {
		return keyIn.(*tables.APIKey), nil
	}
}

func (this *Nasa) updateApiKeyInvokeFre(key *tables.APIKey) error {
	return dao.DefSagaApiDB.ApiDB.UpdateApiKeyInvokeFre(key.ApiKey, key.ApiId, key.UsedNum, this.invokeFreq)
}
