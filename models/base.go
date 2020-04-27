package models

type ApiExtra struct {
	Popularity      int
	Delay           int
	SuccessRate     int
	InvokeFrequency int
}

type ApiDetailInstruction struct {
	DataDesc            string
	DataSource          string
	ApplicationScenario string
}

type QrCodeData struct {
	Action string      `json:"action"`
	Params QrCodeParam `json:"params"`
}

type QrCodeParam struct {
	InvokeConfig InvokeConfig `json:"invokeConfig"`
}

type InvokeConfig struct {
	ContractHash string     `json:"contractHash"`
	Functions    []Function `json:"functions"`
	Payer        string     `json:"payer"`
	GasLimit     uint64     `json:"gasLimit"`
	GasPrice     uint64     `json:"gasPrice"`
}

type Function struct {
	Operation string `json:"operation"`
	Args      []Arg  `json:"args"`
}

type Arg struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ApiKeyInvokeFre struct {
	ApiKey       string `json:"apiKey"`
	OrderId      string `json:"orderId"`
	ApiId        int    `json:"apiId"`
	RequestLimit int    `json:"requestLimit"`
	UsedNum      int32  `json:"usedNum"`
	OntId        string
	InvokeFre    int32
}
