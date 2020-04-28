package nasa

import (
	"errors"
	"fmt"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models"
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
	apiKeyInvokeFreCache *sync.Map //apikey -> ApiKeyInvokeFre
	freqLock             *sync.Mutex
	updateFreq           chan string
}

func NewNasa() *Nasa {
	res := &Nasa{
		apiKeyInvokeFreCache: new(sync.Map),
		freqLock:             new(sync.Mutex),
		updateFreq:           make(chan string, 20),
	}

	go res.UpdateFreqDataBase()
	return res
}

func (this *Nasa) UpdateFreqDataBase() {
	for {
		select {
		case apiKey := <-this.updateFreq:
			keyIn, ok := this.apiKeyInvokeFreCache.Load(apiKey)
			if !ok {
				continue
			}

			key := keyIn.(*models.ApiKeyInvokeFre)
			this.updateApiKeyInvokeFre(key)
		}
	}
}

func (this *Nasa) beforeCheckApiKey(apiKey string, apiId int) (*models.ApiKeyInvokeFre, error) {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyInvokeFre(apiKey)
	if err != nil {
		return nil, err
	}

	if key.UsedNum >= int32(key.RequestLimit) {
		return nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.RequestLimit)
	}
	if key.ApiId != apiId {
		return nil, fmt.Errorf("this apikey: %s can not invoke this api", apiKey)
	}
	key.UsedNum += 1
	key.InvokeFre += 1
	return key, nil
}

func (this *Nasa) Apod(apiKey string) (res []byte, e error) {
	key, err := this.beforeCheckApiKey(apiKey, sagaconfig.APOD)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(apod, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, e = http.DefClient.Get(url)
	if e != nil {
		atomic.AddInt32(&key.UsedNum, -1)
		atomic.AddInt32(&key.InvokeFre, -1)
		return nil, err
	}

	this.updateFreq <- apiKey

	return res, nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) (res []byte, e error) {
	key, err := this.beforeCheckApiKey(apiKey, sagaconfig.FEED)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(feed, startDate, endDate, sagaconfig.DefSagaConfig.NASAAPIKey)
	res, e = http.DefClient.Get(url)
	if e != nil {
		atomic.AddInt32(&key.UsedNum, -1)
		atomic.AddInt32(&key.InvokeFre, -1)
		return nil, err
	}

	this.updateFreq <- apiKey

	return
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

func (this *Nasa) getApiKeyInvokeFre(apiKey string) (*models.ApiKeyInvokeFre, error) {
	keyIn, ok := this.apiKeyInvokeFreCache.Load(apiKey)
	if !ok || keyIn == nil {
		key, err := dao.DefSagaApiDB.ApiDB.QueryApiKeyAndInvokeFreByApiKey(apiKey)
		if err != nil {
			return nil, err
		}
		this.apiKeyInvokeFreCache.Store(apiKey, key)
		return key, nil
	} else {
		return keyIn.(*models.ApiKeyInvokeFre), nil
	}
}

func (this *Nasa) updateApiKeyInvokeFre(key *models.ApiKeyInvokeFre) error {
	return dao.DefSagaApiDB.ApiDB.UpdateApiKeyInvokeFre(key.ApiKey, int(key.UsedNum), key.ApiId, int(key.InvokeFre))
}
