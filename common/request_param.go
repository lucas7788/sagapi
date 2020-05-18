package common

import (
	"github.com/ontio/sagapi/models/tables"
)

type TakeOrderParam struct {
	ProductName      string `json:"productName"`
	OntId            string `json:"ontId"`
	UserName         string `json:"userName"`
	ApiId            uint32 `json:"apiId"`
	SpecificationsId uint32 `json:"specificationsId"`
}

type GenerateTestKeyParam struct {
	ApiId uint32 `json:"apiId"`
}

type AdminGenerateTestKeyParam struct {
	ApiId    uint32 `json:"apiId"`
	ApiState int32  `json:"apiState"`
}

type GetPayQrCode struct {
	OrderId  string `json:"orderId"`
	OntId    string `json:"ontId"`
	UserName string `json:"userName"`
}

type SendTxParam struct {
	Signer    string    `json:"signer"`
	SignedTx  string    `json:"signedTx"`
	ExtraData ExtraData `json:"extraData"`
}

type ExtraData struct {
	Id        string `json:"id"`
	PublicKey string `json:"publickey"`
	OntId     string `json:"ontId"`
}

type GetQrCodeParam struct {
	Id string `json:"id"`
}

type OrderIdParam struct {
	OrderId string `json:"orderId"`
}
type QrCodeIdParam struct {
	QrCodeId string `json:"qrCodeId"`
}

type SearchApiByKey struct {
	Key string `json:"key"`
}

type GetApiByCategoryId struct {
	CategoryId uint32 `json:"categoryId"`
	PageNum    uint32 `json:"pageNum"`
	PageSize   uint32 `json:"pageSize"`
}

type WetherForcastRequest struct {
	TargetDate  int64           `json:"targetDate"`
	Location    tables.Location `json:"location"`
	ToolBoxId   uint32          `json:"toolBoxId"`
	ApiSourceId uint32          `json:"apiSourceId"`
	AlgorithmId uint32          `json:"algorithmId"`
	EnvId       uint32          `json:"envId"`
	Result      string          `json:"result"`
}

type TokenTemplate struct {
	DataIds   string `json:"dataIds"`
	TokenHash string `json:"tokenHash"`
}

type Dtoken struct {
	ResourceId    string        `json:"resourceId"`
	Account       string        `json:"account"`
	TokenTemplate TokenTemplate `json:"tokenTemplate"`
	Number        uint32        `json:"number"`
}

type WetherForcastServiceRequest struct {
	DataUrl       string                 `json:"dataUrl"`
	Header        map[string]interface{} `json:"header"`
	Param         map[string]interface{} `json:"param"`
	RequestMethod string                 `json:"requestMethod"`
	AlgorithmName string                 `json:"algorithmName"`
	Dtoken        Dtoken                 `json:"dtoken"`
}
