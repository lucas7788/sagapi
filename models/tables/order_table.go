package tables

import "github.com/ontio/sagapi/sagaconfig"

type Order struct {
	OrderId          string                 `json:"orderId" db:"OrderId"`
	Title            string                 `json:"title" db:"Title"`
	ProductName      string                 `json:"productName" db:"ProductName"`
	OrderType        string                 `json:"orderType" db:"OrderType"`
	OrderTime        int64                  `json:"orderTime" db:"OrderTime"`
	PayTime          int64                  `json:"payTime" db:"PayTime"`
	OrderStatus      sagaconfig.OrderStatus `json:"orderStatus" db:"OrderStatus"`
	Amount           string                 `json:"amount" db:"Amount"`
	OntId            string                 `json:"ontId" db:"OntId"`
	UserName         string                 `json:"userName" db:"UserName"`
	TxHash           string                 `json:"txHash" db:"TxHash"`
	Price            string                 `json:"price" db:"Price"`
	ApiId            int                    `json:"apiId" db:"ApiId"`
	ApiUrl           string                 `json:"apiUrl" db:"ApiUrl"`
	SpecificationsId int                    `json:"specificationsId" db:"SpecificationsId"`
	Coin             string                 `json:"coin" db:"Coin"`
}

type APIKey struct {
	Id           int                    `json:"id" db:"Id"`
	ApiKey       string                 `json:"apiKey" db:"ApiKey"`
	OrderId      string                 `json:"orderId" db:"OrderId"`
	ApiId        int                    `json:"apiId" db:"ApiId"`
	RequestLimit int                    `json:"requestLimit" db:"RequestLimit"`
	UsedNum      int32                  `json:"usedNum" db:"UsedNum"`
	OntId        string                 `json:"ontId" db:"OntId"`
	OrderStatus  sagaconfig.OrderStatus `json:"orderStatus" db:"OrderStatus"`
}

type QrCode struct {
	Id         int    `json:"id" db:"Id"`
	QrCodeId   string `json:"id" db:"QrCodeId"`
	Ver        string `json:"ver" db:"Ver"`
	OrderId    string `json:"orderId" db:"OrderId"`
	Requester  string `json:"requester" db:"Requester"`
	Signature  string `json:"signature" db:"Signature"`
	Signer     string `json:"signer" db:"Signer"`
	QrCodeData string `json:"data" db:"QrCodeData"`
	Callback   string `json:"callback" db:"Callback"`
	Exp        int64  `json:"exp" db:"Exp"`
	Chain      string `json:"chain" db:"Chain"`
	QrCodeDesc string `json:"desc" db:"QrCodeDesc"`
}
