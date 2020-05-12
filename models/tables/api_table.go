package tables

import (
	"time"
)

// for ApiBasicInfo.ApiState
const (
	API_STATE_INVALID int32 = 0
	API_STATE_BUILTIN int32 = 1
	API_STATE_PUBLISH int32 = 2
	API_STATE_LAST    int32 = 3
)

// for RequestParam.ParamWhere
const (
	URL_PARAM_RESTFUL int32 = 1
	URL_PARAM_QUERY   int32 = 2
	URL_PARAM_BODY    int32 = 3
)

const (
	API_REQUEST_POST string = "POST"
	API_REQUEST_GET  string = "GET"
)

type ApiBasicInfo struct {
	ApiId               uint32    `json:"apiId" db:"ApiId"`
	Coin                string    `json:"coin" db:"Coin"`
	ApiType             string    `json:"type" db:"ApiType"`
	Icon                string    `json:"icon" db:"Icon"`
	Title               string    `json:"title" db:"Title"`
	ApiProvider         string    `json:"provider" db:"ApiProvider"`        //source url. join args can access.
	ApiSagaUrlKey       string    `json:"apiSagaUrlKey" db:"ApiSagaUrlKey"` //sagaurlkey
	ApiUrl              string    `json:"apiUrl" db:"ApiUrl"`               // sagaurl
	Price               string    `json:"price" db:"Price"`
	ApiDesc             string    `json:"description" db:"ApiDesc"`
	ErrorDesc           string    `json:"errorDescription" db:"ErrorDesc"`
	Specifications      uint32    `json:"specifications" db:"Specifications"`
	Popularity          uint32    `json:"popularity" db:"Popularity"`
	Delay               uint32    `json:"delay" db:"Delay"`
	SuccessRate         uint32    `json:"successRate" db:"SuccessRate"`
	InvokeFrequency     uint64    `json:"invokeFrequency" db:"InvokeFrequency"`
	ApiState            int32     `json:"apiState" db:"ApiState"`
	RequestType         string    `json:"requestType" db:"RequestType"`
	Mark                string    `json:"mark" db:"Mark"`
	ResponseParam       string    `json:"responseParam" db:"ResponseParam"`
	ResponseExample     string    `json:"responseExample" db:"ResponseExample"`
	DataDesc            string    `json:"dataDesc" db:"DataDesc"`
	DataSource          string    `json:"dataSource" db:"DataSource"`
	ApplicationScenario string    `json:"applicationScenario" db:"ApplicationScenario"`
	OntId               string    `json:"ontId" db:"OntId"`
	Author              string    `json:"author" db:"Author"`
	CreateTime          time.Time `json:"createTime" db:"CreateTime"`
}

type ApiTag struct {
	Id         uint32    `json:"id" db:"Id"`
	ApiId      uint32    `json:"apiId" db:"ApiId"`
	TagId      uint32    `json:"tagId" db:"TagId"`
	State      byte      `json:"state" db:"State"`
	CreateTime time.Time `json:"createTime" db:"CreateTime"`
}

type Tag struct {
	Id         uint32    `json:"id" db:"Id"`
	Name       string    `json:"name" db:"Name"`
	CategoryId uint32    `json:"categoryId" db:"CategoryId"`
	State      byte      `json:"state" db:"State"`
	CreateTime time.Time `json:"createTime" db:"CreateTime"`
}

type Category struct {
	Id     uint32 `json:"id" db:"Id"`
	NameZh string `json:"nameZh" db:"NameZh"`
	NameEn string `json:"nameEn" db:"NameEn"`
	Icon   string `json:"icon" db:"Icon"`
	State  byte   `json:"state" db:"State"`
}

type Specifications struct {
	Id     uint32 `json:"id" db:"Id"`
	ApiId  uint32 `json:"apiId" db:"ApiId"`
	Price  string `json:"price" db:"Price"`
	Amount uint64 `json:"amount" db:"Amount"`
}

type RequestParam struct {
	Id         uint32 `json:"id" db:"Id"`
	ApiId      uint32 `json:"apiId" db:"ApiId"`
	ParamName  string `json:"paramName" db:"ParamName"`
	Required   bool   `json:"required" db:"Required"`
	ParamWhere int32  `json:"paramWhere" db:"ParamWhere"`
	ParamType  string `json:"paramType" db:"ParamType"`
	Note       string `json:"note" db:"Note"`
	ValueDesc  string `json:"valueDesc" db:"ValueDesc"`
}

type ErrorCode struct {
	Id        uint32 `json:"id" db:"Id"`
	ErrorCode int32  `json:"errorCode" db:"ErrorCode"`
	ErrorDesc string `json:"errorDesc" db:"ErrorDesc"`
}
