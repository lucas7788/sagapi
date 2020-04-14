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
