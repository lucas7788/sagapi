package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"sync/atomic"
)

func HandleDataSourceReqCore(tx *sqlx.Tx, sagaUrlKey string, params []*tables.RequestParam, apiKey string, publishTestOnly bool) ([]byte, error) {
	log.Debugf("HandleDataSourceReqCore : %v", params)
	info, err := dao.DefSagaApiDB.QueryApiBasicInfoBySagaUrlKey(tx, sagaUrlKey)
	if err != nil {
		return nil, err
	}

	referParams, err := dao.DefSagaApiDB.QueryRequestParamByApiId(tx, info.ApiId)
	if err != nil {
		return nil, err
	}

	if len(referParams) != len(params) {
		return nil, fmt.Errorf("params len error. should be %v", referParams)
	}

	var firstQueryArg bool
	var bodyParam []byte
	baseUrl := info.ApiProvider
	firstQueryArg = true
	bodyParamNum := uint32(0)
	for i, p := range referParams {
		log.Debugf("publish param[%d]: %v", i, p)
		if p.ParamName != params[i].ParamName || p.ParamWhere != params[i].ParamWhere {
			return nil, fmt.Errorf("params error. should be %v", referParams)
		}

		switch p.ParamWhere {
		case tables.URL_PARAM_RESTFUL:
			if !firstQueryArg {
				return nil, fmt.Errorf("params error. restful url after query.")
			}
			baseUrl = baseUrl + "/" + params[i].Note
		case tables.URL_PARAM_QUERY:
			if firstQueryArg {
				baseUrl = baseUrl + "?" + params[i].ParamName + "=" + params[i].Note
			} else {
				baseUrl = baseUrl + "&" + params[i].ParamName + "=" + params[i].Note
			}
		case tables.URL_PARAM_BODY:
			if bodyParamNum != 0 {
				return nil, fmt.Errorf("params error. can not pass multi body param.")
			}

			bodyParamNum += 1
			bodyParam = []byte(params[i].Note)
		}
	}

	log.Debugf("baseUrl: %s", baseUrl)

	var key *tables.APIKey
	var apiCounterP *uint64
	if !publishTestOnly {
		key, apiCounterP, err = DefSagaApi.Cache.BeforeCheckApiKey(apiKey, 0)
		if err != nil {
			return nil, err
		}
	}

	switch info.RequestType {
	case tables.API_REQUEST_GET:
		res, err := http.DefClient.Get(baseUrl)
		if err != nil {
			if !publishTestOnly {
				atomic.AddUint64(&key.UsedNum, ^uint64(0))
				atomic.AddUint64(apiCounterP, ^uint64(0))
			}
			return nil, err
		}

		if !publishTestOnly {
			DefSagaApi.Cache.UpdateFreq <- apiKey
		}
		return res, nil
	case tables.API_REQUEST_POST:
		res, err := http.DefClient.Post(baseUrl, bodyParam)
		if err != nil {
			if !publishTestOnly {
				atomic.AddUint64(&key.UsedNum, ^uint64(0))
				atomic.AddUint64(apiCounterP, ^uint64(0))
			}
			return nil, err
		}
		if !publishTestOnly {
			DefSagaApi.Cache.UpdateFreq <- apiKey
		}
		return res, nil
	}

	return nil, nil
}
