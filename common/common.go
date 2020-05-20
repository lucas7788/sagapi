package common

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/ontio/sagapi/utils"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	WETHER_DATA_PROCESS string = "Weather Forecast"
)

const (
	TEST_APIKEY_PREFIX   string = "test_"
	NORMAL_APIKEY_PREFIX string = "apikey_"
	SAGA_URL_PREFIX      string = "sagaurl_"
	ORDER_ID_PREFIX      string = "orderId_"
	QRCODE_ID_PREFIX     string = "qrcodeId_"
)

const (
	UUID_TYPE_RAW          int32 = 1
	UUID_TYPE_TEST_API_KEY int32 = 2
	UUID_TYPE_API_KEY      int32 = 3
	UUID_TYPE_SAGA_URL     int32 = 4
	UUID_TYPE_ORDER_ID     int32 = 5
	UUID_TYPE_QRCODE_ID    int32 = 6
)

type qrCodeDesc struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Price  string `json:"price"`
}

func GenerateUUId(uuidType int32) string {
	u1 := uuid.Must(uuid.NewV4())
	switch uuidType {
	case UUID_TYPE_RAW:
		return u1.String()
	case UUID_TYPE_TEST_API_KEY:
		return TEST_APIKEY_PREFIX + u1.String()
	case UUID_TYPE_SAGA_URL:
		return SAGA_URL_PREFIX + u1.String()
	case UUID_TYPE_API_KEY:
		return NORMAL_APIKEY_PREFIX + u1.String()
	case UUID_TYPE_ORDER_ID:
		return ORDER_ID_PREFIX + u1.String()
	case UUID_TYPE_QRCODE_ID:
		return QRCODE_ID_PREFIX + u1.String()
	}

	return u1.String()
}

func IsTestKey(apiKey string) bool {
	return strings.Contains(apiKey, TEST_APIKEY_PREFIX)
}

func BuildQrCodeResponse(id string) *QrCodeResponse {
	return &QrCodeResponse{
		QrCode: QrCode{
			ONTAuthScanProtocol: sagaconfig.DefSagaConfig.ONTAuthScanProtocol + "/" + id,
		},
		Id: id,
	}
}

func BuildQrCode(orderId, ontid, payer, from, to, value string) (*tables.QrCode, error) {
	return buildQrCode(sagaconfig.DefSagaConfig.NetType, orderId, ontid, payer, from, to, value)
}

func buildQrCode(chain, orderId, ontid, payer, from, to, value string) (*tables.QrCode, error) {
	exp := time.Now().Unix() + sagaconfig.QrCodeExp
	amt := utils.ToIntByPrecise(value, sagaconfig.ONG_DECIMALS)
	data := &models.QrCodeData{
		Action: "signTransaction",
		Params: models.QrCodeParam{
			InvokeConfig: models.InvokeConfig{
				ContractHash: sagaconfig.ONG_CONTRACT_ADDRESS,
				Functions: []models.Function{
					models.Function{
						Operation: "transfer",
						Args: []models.Arg{
							models.Arg{
								Name:  "from",
								Value: "Address:" + from,
							},
							models.Arg{
								Name:  "to",
								Value: "Address:" + to,
							},
							models.Arg{
								Name:  "value",
								Value: amt.Uint64(),
							},
						},
					},
				},
				Payer:    payer,
				GasLimit: 20000,
				GasPrice: 500,
			},
		},
	}
	databs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	id := GenerateUUId(UUID_TYPE_QRCODE_ID)
	sig, err := sagaconfig.DefSagaConfig.OntIdAccount.Sign(databs)
	if err != nil {
		return nil, err
	}

	qrDesc := qrCodeDesc{
		Type:   "invoke neovm type",
		Detail: "transfer",
		Price:  value + "ONG",
	}

	qrDescIn, err := json.Marshal(qrDesc)
	if err != nil {
		return nil, err
	}

	return &tables.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		OrderId:    orderId,
		Requester:  sagaconfig.DefSagaConfig.OntId,
		Signature:  common.ToHexString(sig),
		Signer:     ontid,
		QrCodeData: string(databs),
		Callback:   sagaconfig.DefSagaConfig.QrCodeCallback,
		Exp:        exp,
		Chain:      chain,
		QrCodeDesc: string(qrDescIn),
	}, nil
}

func BuildWetherForcastQrCode(chain, orderId, ontid string, resourceidApi string, resourceidalgenv string, auth_token_templatehex string, use_token_templatehex string, payer, agent string, value string) (*tables.QrCode, error) {
	resourceidApihex := hex.EncodeToString([]byte(resourceidApi))
	resourceidalgenvhex := hex.EncodeToString([]byte(resourceidalgenv))
	resourceids := make([]string, 2, 2)
	resourceids[0] = "ByteArray:" + resourceidApihex
	resourceids[1] = "ByteArray:" + resourceidalgenvhex
	ns := make([]int, 2, 2)
	ns[0] = 1
	ns[1] = 1
	exp := time.Now().Unix() + sagaconfig.QrCodeExp
	data := &models.QrCodeData{
		Action: "signTransaction",
		Params: models.QrCodeParam{
			InvokeConfig: models.InvokeConfig{
				ContractHash: "844648a43e90c641f74255ccc2191b191c4e99a8",
				Functions: []models.Function{
					models.Function{
						Operation: "buyDtokensAndSetAgents",
						Args: []models.Arg{
							models.Arg{
								Name:  "resourceids",
								Value: resourceids,
							},
							models.Arg{
								Name:  "ns",
								Value: ns,
							},
							models.Arg{
								Name:  "use_index",
								Value: uint32(1), // use env/alg
							},
							models.Arg{
								Name:  "authorized_index",
								Value: uint32(0), // auth api
							},
							models.Arg{
								Name:  "authorized_token_template",
								Value: "ByteArray:" + "00" + "20" + auth_token_templatehex, // auth api
							},
							models.Arg{
								Name:  "use_token_template",
								Value: "ByteArray:" + "00" + "20" + use_token_templatehex, // auth api
							},
							models.Arg{
								Name:  "buyer",
								Value: "Address:" + payer,
							},
							models.Arg{
								Name:  "agent",
								Value: "Address:" + agent,
							},
						},
					},
				},
				Payer:    payer,
				GasLimit: 100000,
				GasPrice: 500,
			},
		},
	}
	databs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	log.Debugf("BuildWetherForcastQrCode: %s", string(databs))

	id := GenerateUUId(UUID_TYPE_QRCODE_ID)
	sig, err := sagaconfig.DefSagaConfig.OntIdAccount.Sign(databs)
	if err != nil {
		return nil, err
	}

	qrDesc := qrCodeDesc{
		Type:   "invoke wasm type.",
		Detail: "WetherForcast api process transaction.",
		Price:  value + "ONG",
	}

	qrDescIn, err := json.Marshal(qrDesc)
	if err != nil {
		return nil, err
	}

	return &tables.QrCode{
		Ver:          "1.0.0",
		QrCodeId:     id,
		OrderId:      orderId,
		Requester:    sagaconfig.DefSagaConfig.OntId,
		Signature:    common.ToHexString(sig),
		Signer:       ontid,
		QrCodeData:   string(databs),
		Callback:     sagaconfig.DefSagaConfig.QrCodeCallback,
		Exp:          exp,
		Chain:        chain,
		ContractType: "wasm",
		QrCodeDesc:   string(qrDescIn),
	}, nil
}
