package common

type TakeOrderParam struct {
	ProductName      string `json:"productName"`
	OntId            string `json:"ontId"`
	UserName         string `json:"userName"`
	ApiId            int    `json:"apiId"`
	SpecificationsId int    `json:"specifications"`
}

type GenerateTestKeyParam struct {
	ApiId int `json:"apiId"`
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
	OrderId string
}
type QrCodeIdParam struct {
	QrCodeId string
}
