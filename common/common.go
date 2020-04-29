package common

import (
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

func GenerateUUId() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func IsTestKey(apiKey string) bool {
	return strings.Contains(apiKey, "test")
}

func BuildQrCodeResponse(id string) *QrCodeResponse {
	return &QrCodeResponse{
		QrCode: QrCode{
			ONTAuthScanProtocol: sagaconfig.DefSagaConfig.ONTAuthScanProtocol + "/" + id,
		},
		Id: id,
	}
}

func BuildQrCode(orderId, ontid, payer, from, to, value string) *tables.QrCode {
	return buildQrCode(sagaconfig.DefSagaConfig.NetType, orderId, ontid, payer, from, to, value)
}

func buildQrCode(chain, orderId, ontid, payer, from, to, value string) *tables.QrCode {
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
		//TODO
	}
	log.Errorf("qrdata length: %d", len(databs))
	id := GenerateUUId()
	sig, err := sagaconfig.DefSagaConfig.OntIdAccount.Sign(databs)
	if err != nil {

	}
	return &tables.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		OrderId:    orderId,
		Requester:  sagaconfig.OntId,
		Signature:  common.ToHexString(sig),
		Signer:     ontid,
		QrCodeData: string(databs),
		Callback:   sagaconfig.DefSagaConfig.QrCodeCallback,
		Exp:        exp,
		Chain:      chain,
		QrCodeDesc: "",
	}
}
