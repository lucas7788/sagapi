package tables

import (
	"time"
)

type ApiBasicInfo struct {
	ApiId           int    `json:"apiId" db:"ApiId"`
	Coin            string `json:"coin" db:"Coin"`
	ApiType         string `json:"type" db:"ApiType"`
	Icon            string `json:"icon" db:"Icon"`
	Title           string `json:"title" db:"Title"`
	ApiProvider     string `json:"provider" db:"ApiProvider"`
	ApiUrl          string `json:"apiUrl" db:"ApiUrl"`
	Price           string `json:"price" db:"Price"`
	ApiDesc         string `json:"description" db:"ApiDesc"`
	Specifications  int    `json:"specifications" db:"Specifications"`
	Popularity      int    `json:"popularity" db:"Popularity"`
	Delay           int    `json:"delay" db:"Delay"`
	SuccessRate     int    `json:"successRate" db:"SuccessRate"`
	InvokeFrequency int    `json:"invokeFrequency" db:"InvokeFrequency"`
	CreateTime      string `json:"createTime" db:"CreateTime"`
}

type ApiTag struct {
	Id         int       `json:"id" db:"	Id"`
	ApiId      int       `json:"apiId" db:"ApiId"`
	TagId      int       `json:"tagId" db:"TagId"`
	State      byte      `json:"state" db:"State"`
	CreateTime time.Time `json:"createTime" db:"CreateTime"`
}

type Tag struct {
	Id         int       `json:"id" db:"	Id"`
	Name       string    `json:"name" db:"Name"`
	CategoryId int       `json:"categoryId" db:"CategoryId"`
	State      byte      `json:"state" db:"State"`
	CreateTime time.Time `json:"createTime" db:"CreateTime"`
}

type Category struct {
	Id     int    `json:"id" db:"Id"`
	NameZh string `json:"nameZh" db:"NameZh"`
	NameEn string `json:"nameEn" db:"NameEn"`
	Icon   string `json:"icon" db:"Icon"`
	State  byte   `json:"state" db:"State"`
}

type ApiDetailInfo struct {
	Id                  int    `json:"id" db:"Id"`
	ApiId               int    `json:"apiId" db:"ApiId"`
	RequestType         string `json:"requestType" db:"RequestType"`
	Mark                string `json:"mark" db:"Mark"`
	ResponseParam       string `json:"responseParam" db:"ResponseParam"`
	ResponseExample     string `json:"responseExample" db:"ResponseExample"`
	DataDesc            string `json:"dataDesc" db:"DataDesc"`
	DataSource          string `json:"dataSource" db:"DataSource"`
	ApplicationScenario string `json:"applicationScenario" db:"ApplicationScenario"`
}

type Specifications struct {
	Id              int    `json:"id" db:"Id"`
	ApiDetailInfoId int    `json:"apiDetailInfoId" db:"ApiDetailInfoId"`
	Price           string `json:"price" db:"Price"`
	Amount          int32  `json:"amount" db:"Amount"`
}

type RequestParam struct {
	Id              int    `json:"id" db:"Id"`
	ApiDetailInfoId int    `json:"apiDetailInfoId" db:"ApiDetailInfoId"`
	ParamName       string `json:"paramName" db:"ParamName"`
	Required        bool   `json:"required" db:"Required"`
	ParamType       string `json:"paramType" db:"ParamType"`
	Note            string `json:"note" db:"Note"`
	ValueDesc       string `json:"valueDesc" db:"ValueDesc"`
}

type ErrorCode struct {
	Id              int    `json:"id" db:"Id"`
	ApiDetailInfoId int    `json:"apiDetailInfoId" db:"ApiDetailInfoId"`
	ErrorCode       int    `json:"errorCode" db:"ErrorCode"`
	ErrorDesc       string `json:"errorDesc" db:"ErrorDesc"`
}
