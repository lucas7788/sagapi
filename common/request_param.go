package common

type TakeOrderParam struct {
	ProductName    string `json:"productName"`
	Type           string `json:"type"`
	OntId          string `json:"ontId"`
	UserName       string `json:"userName"`
	Price          string `json:"price"`
	ApiId          int    `json:"apiId"`
	Specifications int    `json:"specifications"`
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
