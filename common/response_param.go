package common

type GetOrderResult struct {
	Result   string `json:"result"`
	UserName string `json:"userName"`
	OntId    string `json:"ontId"`
}

type QrCodeResult struct {
	QrCode QrCode `json:"qrCode"`
	Id     string `json:"id"`
}

type QrCode struct {
	ONTAuthScanProtocol string `json:"ONTAuthScanProtocol"`
}
