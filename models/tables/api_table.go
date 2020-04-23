package tables

type ApiBasicInfo struct {
	ApiId           int    `json:"apiId"`
	Coin            string `json:"coin"` //ONG,ONT
	ApiType         string `json:"type"` //api or data
	Icon            string `json:"icon"`
	Title           string `json:"title"`
	ApiProvider     string `json:"provider"`
	ApiUrl          string `json:"apiUrl"`
	Price           string `json:"price"`
	ApiDesc         string `json:"description"`
	Specifications  int    `json:"specifications"`
	Popularity      int    `json:"popularity"`
	Delay           int    `json:"delay"`
	SuccessRate     int    `json:"successRate"`
	InvokeFrequency int    `json:"invokeFrequency"`
	CreateTime      string `json:"createTime"`
}

type ApiTag struct {
	Id         int  `json:"id"`
	ApiId      int  `json:"apiId"`
	TagId      int  `json:"tagId"`
	State      byte `json:"state"`
	CreateTime int  `json:"createTime"`
}

type Tag struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CategoryId int    `json:"categoryId"`
	State      byte   `json:"state"`
	CreateTime int    `json:"createTime"`
}

type Category struct {
	Id     int    `json:"id"`
	NameZh string `json:"nameZh"`
	NameEn string `json:"nameEn"`
	Icon   string `json:"icon"`
	State  byte   `json:"state"`
}

type ApiDetailInfo struct {
	Id                  int    `json:"id"`
	ApiId               int    `json:"apiId"`
	RequestType         string `json:"requestType"`
	Mark                string `json:"mark"`
	ResponseParam       string `json:"responseParam"`
	ResponseExample     string `json:"responseExample"`
	DataDesc            string `json:"dataDesc"`
	DataSource          string `json:"dataSource"`
	ApplicationScenario string `json:"applicationScenario"`
}

type Specifications struct {
	Id              int    `json:"id"`
	ApiDetailInfoId int    `json:"apiDetailInfoId"`
	Price           string `json:"price"`
	Amount          int    `json:"amount"`
}

type RequestParam struct {
	ApiDetailInfoId int    `json:"apiDetailInfoId"`
	ParamName       string `json:"paramName"`
	Required        bool   `json:"required"`
	ParamType       string `json:"paramType"`
	Note            string `json:"note"`
	ValueDesc       string `json:"valueDesc"`
}

type ErrorCode struct {
	ApiDetailInfoId int    `json:"apiDetailInfoId"`
	ErrorCode       int    `json:"errorCode"`
	ErrorDesc       string `json:"errorDesc"`
}
