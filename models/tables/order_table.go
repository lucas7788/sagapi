package tables

import "github.com/ontio/sagapi/sagaconfig"

type Order struct {
	OrderId          string                 `json:"orderId"`
	Title            string                 `json:"title"`
	ProductName      string                 `json:"productName"`
	OrderType        string                 `json:"orderType"`
	OrderTime        int64                  `json:"orderTime"`
	PayTime          int64                  `json:"payTime"`
	OrderStatus      sagaconfig.OrderStatus `json:"orderStatus"`
	Amount           string                 `json:"amount"`
	OntId            string                 `json:"ontId"`
	UserName         string                 `json:"userName"`
	TxHash           string                 `json:"txHash"`
	Price            string                 `json:"price"`
	ApiId            int                    `json:"apiId"`
	SpecificationsId int                    `json:"specificationsId"`
	Coin             string                 `json:"coin"`
}

type APIKey struct {
	Id           int                    `json:"id"`
	ApiKey       string                 `json:"apiKey"`
	OrderId      string                 `json:"orderId"`
	ApiId        int                    `json:"apiId"`
	RequestLimit int                    `json:"requestLimit"`
	UsedNum      int                    `json:"usedNum"`
	OntId        string                 `json:"ontId"`
	OrderStatus  sagaconfig.OrderStatus `json:"orderStatus"`
}

type QrCode struct {
	QrCodeId   string `json:"id"`
	Ver        string `json:"ver"`
	OrderId    string `json:"orderId"`
	Requester  string `json:"requester"`
	Signature  string `json:"signature"`
	Signer     string `json:"signer"`
	QrCodeData string `json:"data"`
	Callback   string `json:"callback"`
	Exp        int64  `json:"exp"`
	Chain      string `json:"chain"`
	QrCodeDesc string `json:"desc"`
}
