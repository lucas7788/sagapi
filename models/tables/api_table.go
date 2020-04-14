package tables

import (
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models"
)

type Order struct {
	Id          uint `gorm:"primary_key"`
	OrderId     string
	ProductName string
	Type        string
	OrderTime   int64
	PaiedTime   int64
	OrderStatus config.OrderStatus
	Amount      string
	OntId       string
	UserName    string
	TxHash      string
	Price       string
	ApiId       uint
	ApiKey      string
}

type ApiTestRecord struct {
	ID         uint `gorm:"primary_key"`
	OntId      string
	UserName   string
	ApiId      int
	TestResult int //0 test failed, 1 test success
}

type ApiBasicInfo struct {
	ApiId          uint `gorm:"primary_key"`
	ApiLogo        string
	ApiName        string
	ApiProvider    string
	ApiUrl         string
	ApiPrice       string
	ApiDesc        string
	Specifications int
	ApiExtra       models.ApiExtra
}

type ApiDetailInfo struct {
	ID                   uint `gorm:"primary_key"`
	Mark                 string
	RequestParam         string
	ResponseParam        string
	ResponseExample      string
	ParamErrorCode       string
	APIDetailInstruction models.ApiDetailInstruction
}

type APIKey struct {
	ID       uint `gorm:"primary_key"`
	ApiKey   string
	ApiId    int
	Limit    int
	UsedNum  int
	OntId    string
	UserName string
}

type QrCode struct {
	QrCodeId  uint               `gorm:"primary_key"`
	Ver       string             `json:"ver"`
	Id        string             `json:"id"`
	OrderId   string             `json:"orderId"`
	Requester string             `json:"requester"`
	Signature string             `json:"signature"`
	Signer    string             `json:"signer"`
	Data      *models.QrCodeData `json:"data"`
	Callback  string             `json:"callback"`
	Exp       int64              `json:"exp"`
	Chain     string             `json:"chain"`
	Desc      string             `json:"desc"`
}
