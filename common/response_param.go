package common

import (
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
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
	ApiId               int                      `json:"apiId"`
	Mark                string                   `json:"mark"`
	ResponseParam       string                   `json:"responseParam"`
	ResponseType        string                   `json:"responseType"`
	ResponseExample     string                   `json:"responseExample"`
	DataDesc            string                   `json:"dataDesc"`
	DataSource          string                   `json:"dataSource"`
	ApplicationScenario string                   `json:"applicationScenario"`
	RequestParams       []*tables.RequestParam   `json:"requestParams"`
	ErrorCodes          []*tables.ErrorCode      `json:"errorCodes"`
	Specifications      []*tables.Specifications `json:"specifications"`
	ApiBasicInfo        *tables.ApiBasicInfo     `json:"apiBasicInfo"`
}

type OrderResult struct {
	Title        string                 `json:"title"`
	Total        int                    `json:"total"`
	OrderId      string                 `json:"orderId"`
	Amount       string                 `json:"amount`
	CreateTime   int64                  `json:"createTime"`
	TxHash       string                 `json:"txHash"`
	ApiId        int                    `json:"apiId"`
	ApiUrl       string                 `json:"apiUrl"`
	RequestLimit int32                    `json:"requestLimit"`
	UsedNum      int32                    `json:"usedNum"`
	Status       sagaconfig.OrderStatus `json:"status"`
	ApiKey       string                 `json:"apiKey"`
	Price        string                 `json:"price"`
	Coin         string                 `json:"coin"`
}
