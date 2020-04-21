package common

import "github.com/ontio/sagapi/models/tables"

type GetOrderResponse struct {
	Result   string `json:"result"`
	UserName string `json:"userName"`
	OntId    string `json:"ontId"`
}

type QrCodeResponse struct {
	QrCode QrCode `json:"qrCode"`
	Id     string `json:"id"`
}

type QrCode struct {
	ONTAuthScanProtocol string `json:"ONTAuthScanProtocol"`
}

type ApiDetailResponse struct {
	ApiId               int
	Mark                string
	ResponseParam       string
	ResponseType        string
	ResponseExample     string
	DataDesc            string
	DataSource          string
	ApplicationScenario string
	RequestParams       []*tables.RequestParam
	ErrorCodes          []*tables.ErrorCode
	Specifications      []*tables.Specifications
	ApiBasicInfo        *tables.ApiBasicInfo
}
