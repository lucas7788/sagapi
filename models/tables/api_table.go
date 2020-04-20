package tables

type ApiBasicInfo struct {
	ApiId           int
	Coin            string //ONG,ONT
	ApiType         string //api or data
	Icon            string
	Title           string
	ApiProvider     string
	ApiUrl          string
	Price           string
	ApiDesc         string
	Specifications  int
	Popularity      int
	Delay           int
	SuccessRate     int
	InvokeFrequency int
	CreateTime      int
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
	Mark                string
	ResponseParam       string
	ResponseExample     string
	DataDesc            string
	DataSource          string
	ApplicationScenario string
}

type RequestParam struct {
	ApiDetailInfoId int
	ParamName       string
	Required        int8
	ParamType       string
	Note            string
}

type ErrorCode struct {
	ApiDetailInfoId int
	ErrorCode       int
	ErrorDesc       string
}
