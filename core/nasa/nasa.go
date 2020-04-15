package nasa

import (
	"fmt"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"sync"
)

var (
	apod = "https://api.nasa.gov/planetary/apod?api_key=%s"
	feed = "https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s"
)

type Nasa struct {
	apiKeyCache *sync.Map //apikey -> ApiKey
}

func NewNasa() *Nasa {
	return &Nasa{
		apiKeyCache: new(sync.Map),
	}
}

func (this *Nasa) beforeCheckApiKey(apiKey string) (*tables.APIKey, error) {
	key, err := this.getApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	if key.UsedNum >= key.Limit {
		return nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.Limit)
	}
	return key, nil
}

func (this *Nasa) Apod(apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(apod, config.NASA_API_KEY)
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
	return res, nil
}

func (this *Nasa) Feed(startDate, endDate string, apiKey string) ([]byte, error) {
	key, err := this.beforeCheckApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(feed, startDate, endDate, config.NASA_API_KEY)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//TODO
	err = this.updateApiKey(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *Nasa) getApiKey(apiKey string) (*tables.APIKey, error) {
	keyIn, ok := this.apiKeyCache.Load(apiKey)
	var key *tables.APIKey
	if !ok || keyIn == nil {
		var err error
		key, err = dao.DefDB.QueryApiKeyInfo(apiKey)
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
	return dao.DefDB.UpdateApiKey(key)
}
