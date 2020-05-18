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
	ApiId               uint32                   `json:"apiId"`
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
	Total        uint32                 `json:"total"`
	OrderId      string                 `json:"orderId"`
	Amount       string                 `json:"amount"`
	CreateTime   int64                  `json:"createTime"`
	TxHash       string                 `json:"txHash"`
	ApiId        uint32                 `json:"apiId"`
	ApiUrl       string                 `json:"apiUrl"`
	RequestLimit uint64                 `json:"requestLimit"`
	UsedNum      uint64                 `json:"usedNum"`
	Status       sagaconfig.OrderStatus `json:"status"`
	ApiKey       string                 `json:"apiKey"`
	Price        string                 `json:"price"`
	Coin         string                 `json:"coin"`
}

type DataProcessOrderResult struct {
	Title     string                 `json:"title"`
	OrderId   string                 `json:"orderId"`
	OrderTime int64                  `json:"orderTime"`
	TxHash    string                 `json:"txHash"`
	ApiId     uint32                 `json:"apiId"`
	Status    sagaconfig.OrderStatus `json:"status"`
	Price     string                 `json:"price"`
	Coin      string                 `json:"coin"`
	OrderKind uint32                 `json:"orderKind"`
	Request   string                 `json:"request"`
	Result    string                 `json:"result"`
}

type OrderDetailResponse struct {
	Res interface{}
}

type WetherOrderDetail struct {
	TargetDate int64                `json:"targetDate"`
	Location   *tables.Location     `json:"location"`
	ToolBox    *tables.ToolBox      `json:"toolBoxId"`
	ApiSource  *tables.ApiBasicInfo `json:"apiSourceId"`
	Algorithm  *tables.Algorithm    `json:"algorithmId"`
	Env        *tables.Env          `json:"envId"`
	Result     string               `json:"result"`
}
