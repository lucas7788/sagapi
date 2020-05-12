package publish

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
	"strconv"
	"strings"
)

type UrlParams struct {
	Name  string
	Type  int32
	Index uint32
}

func PublishAPIHandle(c *gin.Context) {
	param := &core.PublishAPI{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	err = core.PublishAPIHandleCore(param)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func VerifyAPIHandle(c *gin.Context) {
	res, err := common.ParseGetParamByParamName(c, "apiId", "sagaUrlKey")
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if len(res) != 2 {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("need pass apiId and sagaUrlKey.")))
		return
	}
	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	sagaUrlKey := res[1]
	if sagaUrlKey == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("sagaUrlKey can not empty.")))
		return
	}
	err = dao.DefSagaApiDB.ApiBasicUpateApiState(nil, tables.API_STATE_BUILTIN, uint32(apiId), sagaUrlKey)
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ApiBasicUpateApiState error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
}

func GetApiDetailByApiIdApiState(c *gin.Context) {
	res, err := common.ParseGetParamByParamName(c, "apiId", "apiState")
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	apiState, err := strconv.Atoi(res[1])
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	info, err := core.DefSagaApi.QueryApiDetailInfoByApiId(uint32(apiId), int32(apiState))
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] QueryApiDetailInfoByApiId error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(info))
}

func GetALLPublishPage(c *gin.Context) {
	arr, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	pageNum, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	infos, err := dao.DefSagaApiDB.QueryApiBasicInfoByPage(uint32(pageNum), uint32(pageSize), tables.API_STATE_PUBLISH)
	if err != nil {
		log.Errorf("[GetALLPublishPage] QueryApiBasicInfoByPage error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}

func ParseUrl(url string) ([]UrlParams, error) {
	params := make([]UrlParams, 0)
	key := url
	count := uint32(0)
	var i, j, k int
	var queryArgHandled bool
	for {
		i = strings.IndexAny(key, "{")
		if i != -1 {
			if queryArgHandled {
				return nil, errors.New("url errror")
			}
			key = key[i:]
			i = strings.IndexAny(key, "{")
			k = strings.IndexAny(key, "}")
			j = strings.IndexAny(key, "/")
			if k == -1 || (j != -1 && k+1 != j) {
				return nil, errors.New("url error")
			}
			p := key[i+1 : k]
			if i+1 == k {
				return nil, errors.New("url error")
			}
			params = append(params, UrlParams{
				Name: p,
				Type: tables.URL_PARAM_RESTFUL,
			})
			count += 1
			key = key[k:]
		} else {
			if !queryArgHandled {
				i = strings.IndexAny(key, "?")
				if i == -1 {
					break
				}

				k = strings.IndexAny(key, "&")
				if k == -1 {
					k = len(key)
				}

				p := key[i+1 : k]
				if i+1 == k {
					return nil, errors.New("url error")
				}
				params = append(params, UrlParams{
					Name: p,
					Type: tables.URL_PARAM_QUERY,
				})
				count += 1
				key = key[k:]
			}
			queryArgHandled = true

			i = strings.IndexAny(key, "&")
			if i == -1 {
				break
			}

			k = strings.IndexAny(key[i+1:], "&")
			if k == -1 {
				k = len(key)
			} else {
				k += 1
			}

			if i+1 == k {
				return nil, errors.New("url error")
			}
			p := key[i+1 : k]
			params = append(params, UrlParams{
				Name: p,
				Type: tables.URL_PARAM_QUERY,
			})
			count += 1
			key = key[k:]
		}
	}

	return params, nil
}
