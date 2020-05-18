package tables

import "github.com/ontio/sagapi/sagaconfig"

const (
	ORDER_KIND_API                 = 1
	ORDER_KIND_DATA_PROCESS_WETHER = 2
	ORDER_KIND_DATA_LAST           = 3
)

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
	ApiId            uint32                 `json:"apiId" db:"ApiId"`
	ApiUrl           string                 `json:"apiUrl" db:"ApiUrl"`
	SpecificationsId uint32                 `json:"specificationsId" db:"SpecificationsId"`
	Coin             string                 `json:"coin" db:"Coin"`
	OrderKind        uint32                 `json:"orderKind" db:"OrderKind"`
	Request          string                 `json:"request" db:"Request"` // this fill
	Result           string                 `json:"result" db:"Result"`   // this fill
}

type APIKey struct {
	Id           uint32 `json:"id" db:"Id"`
	ApiKey       string `json:"apiKey" db:"ApiKey"`
	OrderId      string `json:"orderId" db:"OrderId"`
	ApiId        uint32 `json:"apiId" db:"ApiId"`
	RequestLimit uint64 `json:"requestLimit" db:"RequestLimit"`
	UsedNum      uint64 `json:"usedNum" db:"UsedNum"`
	OntId        string `json:"ontId" db:"OntId"`
}

type QrCode struct {
	QrCodeId     string `json:"id" db:"QrCodeId"`
	Ver          string `json:"ver" db:"Ver"`
	OrderId      string `json:"orderId" db:"OrderId"`
	Requester    string `json:"requester" db:"Requester"`
	Signature    string `json:"signature" db:"Signature"`
	Signer       string `json:"signer" db:"Signer"`
	QrCodeData   string `json:"data" db:"QrCodeData"`
	Callback     string `json:"callback" db:"Callback"`
	Exp          int64  `json:"exp" db:"Exp"`
	Chain        string `json:"chain" db:"Chain"`
	QrCodeDesc   string `json:"desc" db:"QrCodeDesc"`
	ContractType string `json:"contractType" db:"ContractType"`
}
