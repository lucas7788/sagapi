package publish

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
	"strings"
)

const (
	sagaurlPrefix string = "sagaurl_"
)

// apikey in here.
type PublishAPIParam struct {
	ParamName  string `json:"paramName"`
	ParamWhere int32  `json:"paramWhere"`
	ParamType  string `json:"paramType"`
	Required   bool   `json:"required"`
	Note       string `json:"note"`
	ValueDesc  string `json:"valueDesc"`
}

type PublishErrorCode struct {
	Code int32  `json:"code"`
	Desc string `json:"description"`
}

type PublishAPI struct {
	Name            string                  `json:"name"`
	Desc            string                  `json:"description"`
	RequestType     string                  `json:"requestType"`
	ApiProvider     string                  `json:"apiProvider"`
	Tag             []tables.Tag            `json:"tags"`
	DataSource      string                  `json:"dataSource"`
	ResponseExample string                  `json:"responseExample"`
	Error           []PublishErrorCode      `json:"errorCode"`
	Params          []tables.RequestParam   `json:"params"`
	Spec            []tables.Specifications `json:"specifications"`
}

type UrlParams struct {
	Name  string
	Type  int32
	Index uint32
}

func PublishAPIHandle(c *gin.Context) {
	param := &PublishAPI{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	// handle error
	errorDesc, err := json.Marshal(param.Error)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	tags := make([]*tables.Tag, 0)
	for _, tag := range param.Tag {
		t, err := dao.DefSagaApiDB.QueryTag(tag.CategoryId, tag.Name)
		if err != nil {
			common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		}
		tags = append(tags, t)
	}

	apibasic := &tables.ApiBasicInfo{
		Coin:          "ONG",
		Title:         param.Name,
		ApiProvider:   param.ApiProvider,
		ApiSagaUrlKey: sagaurlPrefix + common2.GenerateUUId(),
		ApiDesc:       param.Desc,
		ApiState:      tables.API_STATE_PUBLISH,
		ErrorDesc:     string(errorDesc),
	}

	err = dao.DefSagaApiDB.ApiDB.InsertApiBasicInfo([]*tables.ApiBasicInfo{apibasic})
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.SQL_ERROR, err))
		return
	}
	info, err := dao.DefSagaApiDB.ApiDB.QueryApiBasicInfoBySagaUrlKey(apibasic.ApiSagaUrlKey)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.SQL_ERROR, err))
		return
	}
	apidetail := &tables.ApiDetailInfo{
		ApiId:           info.ApiId,
		RequestType:     param.RequestType,
		ResponseExample: param.ResponseExample,
		DataSource:      param.DataSource,
	}

	_ = dao.DefSagaApiDB.ApiDB.InsertApiDetailInfo(apidetail)
	apidetail2, err = dao.DefSagaApiDB.ApiDB.QueryApiDetailInfoByApiId(apidetail.ApiId)
	// to do. should atomic. use event.

	// tag handle
	for _, apiTag := range tags {
		tag := &tables.ApiTag{
			ApiId: info.ApiId,
			TagId: apiTag.Id,
			State: byte(1),
		}
		_ = dao.DefSagaApiDB.InsertApiTag(tag)
	}
	// spec handle.

	// handle param

	common.WriteResponse(c, common.ResponseSuccess(nil))
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
