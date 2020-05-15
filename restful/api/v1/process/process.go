package process

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
	"strconv"
)

type AlgorithmObj struct {
	Algorithm *tables.Algorithm `json:"algorithm"`
	Env       []*tables.Env     `json:"env"`
}

type ApiSourceObj struct {
	ApiSource  *tables.ApiBasicInfo `json:"apiSource"`
	Algorithms []*AlgorithmObj      `json:"algorithm"`
}

type WetherForcastResponse struct {
	ToolBox      *tables.ToolBox
	ApiSourceObj []*ApiSourceObj        `json:"apiSourceObj"`
	ApiALL       []*tables.ApiBasicInfo `json:"apiALL"`
	AlgorithmALL []*tables.Algorithm    `json:"algorithmALL"`
	EnvAll       []*tables.Env          `json:"envAll"`
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

func GetAllToolBox(c *gin.Context) {
	res, err := dao.DefSagaApiDB.QueryToolBoxAll(nil)
	if err != nil {
		log.Errorf("[GetAllToolBox]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetWetherForcastInfo(c *gin.Context) {
	toolid := c.Param("toolid")
	if toolid == "" {
		log.Errorf("[GetWetherForcastInfo]: %s", errors.New("toolid can not empty."))
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("toolid can not empty.")))
		return
	}
	toolboxid, err := strconv.Atoi(toolid)
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	toolbox, err := dao.DefSagaApiDB.QueryToolBoxById(nil, uint32(toolboxid))
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	algAll := make([]*tables.Algorithm, 0)
	envAll := make([]*tables.Env, 0)

	apis, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiTypeKind(nil, toolbox.Title, tables.API_KIND_DATA_PROCESS, tables.API_STATE_BUILTIN)
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	apisourceresp := make([]*ApiSourceObj, 0)
	for _, api := range apis {
		apialgs, err := dao.DefSagaApiDB.QueryApiAlgorithmsByApiId(nil, api.ApiId)
		if err != nil {
			log.Errorf("[GetWetherForcastInfo]: %s", err)
			common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
			return
		}
		algsresp := make([]*AlgorithmObj, 0)
		for _, apialg := range apialgs {
			alg, err := dao.DefSagaApiDB.QueryAlgorithmById(nil, apialg.AlgorithmId)
			if err != nil {
				log.Errorf("[GetWetherForcastInfo]: %s", err)
				common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
				return
			}

			algenvs, err := dao.DefSagaApiDB.QueryAlgorithmEnvByAlgorithmId(nil, alg.Id)
			if err != nil {
				log.Errorf("[GetWetherForcastInfo]: %s", err)
				common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
				return
			}
			envsresp := make([]*tables.Env, 0)
			for _, env := range algenvs {
				env, err := dao.DefSagaApiDB.QueryEnvById(nil, env.Id)
				if err != nil {
					log.Errorf("[GetWetherForcastInfo]: %s", err)
					common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
					return
				}
				envsresp = append(envsresp, env)
				envAll = append(envAll, env)
			}
			algsresp = append(algsresp, &AlgorithmObj{
				Algorithm: alg,
				Env:       envsresp,
			})
			algAll = append(algAll, alg)
		}
		apisourceresp = append(apisourceresp, &ApiSourceObj{
			ApiSource:  api,
			Algorithms: algsresp,
		})
	}

	res := WetherForcastResponse{
		ToolBox:      toolbox,
		ApiSourceObj: apisourceresp,
		ApiALL:       apis,
		EnvAll:       envAll,
		AlgorithmALL: algAll,
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}
