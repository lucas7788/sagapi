package freq

import (
	"fmt"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"sync"
	"sync/atomic"
)

type DBCache struct {
	UpdateFreq   chan string
	apiKeyCache  *sync.Map //apikey -> ApiKey
	freqLock     *sync.Mutex
	apiFreqCache *sync.Map //ApiID -> int32
}

func NewDBCache() *DBCache {
	res := &DBCache{
		apiKeyCache:  new(sync.Map),
		apiFreqCache: new(sync.Map),
		freqLock:     new(sync.Mutex),
		UpdateFreq:   make(chan string, 20),
	}

	go res.UpdateFreqDataBase()
	return res
}

func (this *DBCache) UpdateFreqDataBase() {
	for {
		select {
		case apiKey := <-this.UpdateFreq:
			keyIn, ok := this.apiKeyCache.Load(apiKey)
			if !ok {
				log.Debugf("apikey cache not exist")
				continue
			}

			key := keyIn.(*tables.APIKey)
			apiId := key.ApiId

			apiCounterP, ok := this.apiFreqCache.Load(apiId)
			if !ok {
				log.Debugf("apicounter cache not exist")
				continue
			}

			counter := atomic.LoadUint64(apiCounterP.(*uint64))
			log.Debugf("UpdateFreqDataBase: apiId: %d, counterP: %v,counter: %d, usedNum:%d, apiKey: %s", apiId, apiCounterP, counter, key.UsedNum, apiKey)
			this.updateApiKeyInvokeFre(key, counter)
		}
	}
}

func (this *DBCache) getApiIdFreqCounter(ApiId uint32) (*uint64, error) {
	apiCounterP, ok := this.apiFreqCache.Load(ApiId)
	if !ok || apiCounterP == nil {
		freq, err := dao.DefSagaApiDB.QueryInvokeFreByApiId(nil, ApiId)
		if err != nil {
			return nil, err
		}
		this.apiFreqCache.Store(ApiId, &freq)
		return &freq, nil
	} else {
		return apiCounterP.(*uint64), nil
	}
}

// better check test api and api key id conflict.
func (this *DBCache) BeforeCheckApiKey(apiKey string, apiId uint32) (*tables.APIKey, *uint64, error) {
	this.freqLock.Lock()
	defer this.freqLock.Unlock()
	key, err := this.getApiKeyCache(apiKey)
	if err != nil {
		return nil, nil, err
	}

	if apiId != 0 && key.ApiId != apiId {
		return nil, nil, fmt.Errorf("this apikey: %s can not invoke this api", apiKey)
	}

	apiId = key.ApiId

	apiCounterP, err := this.getApiIdFreqCounter(apiId)
	if err != nil {
		return nil, nil, err
	}

	if key.UsedNum >= key.RequestLimit {
		return nil, nil, fmt.Errorf("apikey: %s, useNum: %d, limit:%d", apiKey, key.UsedNum, key.RequestLimit)
	}

	key.UsedNum += 1
	if !common.IsTestKey(apiKey) {
		atomic.AddUint64(apiCounterP, 1)
	}

	return key, apiCounterP, nil
}

func (this *DBCache) getApiKeyCache(apiKey string) (*tables.APIKey, error) {
	keyIn, ok := this.apiKeyCache.Load(apiKey)
	if !ok || keyIn == nil {
		key, err := dao.DefSagaApiDB.QueryApiKeyByApiKey(nil, apiKey)
		if err != nil {
			return nil, err
		}
		this.apiKeyCache.Store(apiKey, key)
		return key, nil
	} else {
		return keyIn.(*tables.APIKey), nil
	}
}

func (this *DBCache) updateApiKeyInvokeFre(key *tables.APIKey, freqCounter uint64) error {
	return dao.DefSagaApiDB.UpdateApiKeyInvokeFre(nil, key.ApiKey, key.ApiId, key.UsedNum, freqCounter)
}
