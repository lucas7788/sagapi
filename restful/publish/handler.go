package publish

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ontio/sagapi/core"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/restful/api/common"
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
}

func GetALLPublishPage(c *gin.Context) {
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
