package tables

import "github.com/ontio/sagapi/config"

type Order struct {
	OrderId          string
	ProductName      string
	OrderType        string
	OrderTime        int64
	PayTime          int64
	OrderStatus      config.OrderStatus
	Amount           string
	OntId            string
	UserName         string
	TxHash           string
	Price            string
	ApiId            int
	SpecificationsId int
	Coin             string
}

type APIKey struct {
	Id           int
	ApiKey       string
	OrderId      string
	ApiId        int
	RequestLimit int
	UsedNum      int
	OntId        string
	OrderStatus  config.OrderStatus
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
