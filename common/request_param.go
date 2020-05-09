package common

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
