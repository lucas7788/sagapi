package tables

import (
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models"
)

type Order struct {
	Id             uint `gorm:"primary_key"`
	OrderId        string
	ProductName    string
	Type           string
	OrderTime      int64
	PaiedTime      int64
	OrderStatus    config.OrderStatus
	Amount         string
	OntId          string
	UserName       string
	TxHash         string
	Price          string
	ApiId          uint
	ApiKey         string
	Specifications int
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
	Id                   uint `gorm:"primary_key"`
	ApiId                uint
	Mark                 string
	RequestParamId       uint
	ResponseParam        string
	ResponseExample      string
	ErrorCodeId          uint
	APIDetailInstruction models.ApiDetailInstruction
}

type RequestParam struct {
	Id              uint `gorm:"primary_key"`
	ApiDetailInfoId uint
	Name            string
	Required        bool
	Type            string
	Note            string
}

type ErrorCode struct {
	Id              uint `gorm:"primary_key"`
	ApiDetailInfoId uint
	ErrorCode       int
	ErrorDesc       string
}

type APIKey struct {
	Id       uint `gorm:"primary_key"`
	ApiKey   string
	ApiId    uint
	Limit    int
	UsedNum  int
	OntId    string
}

type QrCode struct {
	Index     uint               `gorm:"primary_key"`
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
