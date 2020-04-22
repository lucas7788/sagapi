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
	Id         int
	ApiId      int
	TagId      int
	State      byte
	CreateTime int
}

type Tag struct {
	Id         int
	Name       string
	CategoryId int
	State      byte
	CreateTime int
}
type Category struct {
	Id     int
	NameZh string
	NameEn string
	Icon   string
	State  byte
}

type ApiDetailInfo struct {
	Id                  int
	ApiId               int
	RequestType         string
	Mark                string
	ResponseParam       string
	ResponseExample     string
	DataDesc            string
	DataSource          string
	ApplicationScenario string
}

type Specifications struct {
	Id              int
	ApiDetailInfoId int
	Price           string
	Amount          int
}

type RequestParam struct {
	ApiDetailInfoId int
	ParamName       string
	Required        bool
	ParamType       string
	ValueDesc       string
}

type ErrorCode struct {
	ApiDetailInfoId int
	ErrorCode       int
	ErrorDesc       string
}
