package tables

type ApiBasicInfo struct {
	ApiId           int
	ApiLogo         string
	ApiName         string
	ApiProvider     string
	ApiUrl          string
	ApiPrice        string
	ApiDesc         string
	Specifications  int
	Popularity      int
	Delay           int
	SuccessRate     int
	InvokeFrequency int
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
	Id              int
	ApiDetailInfoId int
	ParamName       string
	Required        int8
	ParamType       string
	Note            string
}

type ErrorCode struct {
	Id              int
	ApiDetailInfoId int
	ErrorCode       int
	ErrorDesc       string
}
