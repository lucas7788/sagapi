package common

import (
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
)

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

type OrderResult struct {
	Title        string
	Total        int
	OrderId      string
	Amount       string
	CreateTime   int64
	TxHash       string
	ApiId        int
	RequestLimit int
	UsedNum      int
	Status       config.OrderStatus
	ApiKey       string
	Price        string
	Coin         string
}
