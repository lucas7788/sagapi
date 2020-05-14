package process

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
)

type City struct {
	Id      uint32 `json:"id" db:"Id"`
	Country string `json:"Country" db:"Country"`
	City    string `json:"City" db:"City"`
	Lat     string `json:"Lat" db:"Lat"`
	Lng     string `json:"Lng" db:"Lng"`
}

type AlgorithmObj struct {
	Algorithm *tables.Algorithm `json:"algorithm"`
	Env       []*tables.Env     `json:"env"`
}

type ApiSourceObj struct {
	ApiSource  *tables.ApiBasicInfo `json:"apiSource"`
	Algorithms []*AlgorithmObj      `json:"algorithm"`
}

type WetherForcastResponse struct {
	ApiSourceObj []*ApiSourceObj        `json:"apiSourceObj"`
	ApiALL       []*tables.ApiBasicInfo `json:"apiALL"`
	AlgorithmALL []*tables.Algorithm    `json:"algorithmALL"`
	EnvAll       []*tables.Env          `json:"envAll"`
}

type WetherForcastRequest struct {
}

func GetLocation(c *gin.Context) {
	country := c.Param("country")

	res, err := dao.DefSagaApiDB.QueryLocationOfCountryCity(nil, country)
	if err != nil {
		log.Errorf("[GetLocation]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetWetherForcastInfo(c *gin.Context) {
	apiType := c.Param("preditype")
	if apiType == "" {
		log.Errorf("[GetWetherForcastInfo]: %s", errors.New("preditype can not empty."))
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("apiType can not empty.")))
		return
	}
	algsresp := make([]*AlgorithmObj, 0)
	apisourceresp := make([]*ApiSourceObj, 0)
	envsresp := make([]*tables.Env, 0)

	apis, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiTypeKind(nil, apiType, tables.API_KIND_DATA_PROCESS, tables.API_STATE_BUILTIN)
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	for _, api := range apis {
		apialgs, err := dao.DefSagaApiDB.QueryApiAlgorithmsByApiId(nil, api.ApiId)
		if err != nil {
			log.Errorf("[GetWetherForcastInfo]: %s", err)
			common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
			return
		}
		for _, apialg := range apialgs {
			algs, err := dao.DefSagaApiDB.QueryAlgorithmsById(nil, apialg.AlgorithmId)
			if err != nil {
				log.Errorf("[GetWetherForcastInfo]: %s", err)
				common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
				return
			}

			for _, alg := range algs {
				algenvs, err := dao.DefSagaApiDB.QueryAlgorithmEnvByAlgorithmId(nil, alg.Id)
				if err != nil {
					log.Errorf("[GetWetherForcastInfo]: %s", err)
					common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
					return
				}
				for _, env := range algenvs {
					envs, err := dao.DefSagaApiDB.QueryEnvsById(nil, env.Id)
					if err != nil {
						log.Errorf("[GetWetherForcastInfo]: %s", err)
						common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
						return
					}
					envsresp = envs
				}
				algsresp = append(algsresp, &AlgorithmObj{
					Algorithm: alg,
					Env:       envsresp,
				})
			}
		}
		apisourceresp = append(apisourceresp, &ApiSourceObj{
			ApiSource:  api,
			Algorithms: algsresp,
		})
	}
	envxs, err0 := dao.DefSagaApiDB.QueryEnvsById(nil, 0)
	algxs, err1 := dao.DefSagaApiDB.QueryAlgorithmsById(nil, 0)
	if err0 != nil || err1 != nil {
		log.Errorf("[GetWetherForcastInfo]: %s, %s", err0, err1)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err1))
		return
	}
	res := WetherForcastResponse{
		ApiSourceObj: apisourceresp,
		ApiALL:       apis,
		EnvAll:       envxs,
		AlgorithmALL: algxs,
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}
