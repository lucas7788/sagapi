package common

import (
	"encoding/json"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"github.com/satori/go.uuid"
	"time"
)

func GenerateUUId() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func BuildQrCodeResult(id string) *QrCodeResponse {
	return &QrCodeResponse{
		QrCode: QrCode{
			ONTAuthScanProtocol: config.DefConfig.ONTAuthScanProtocol,
		},
		Id: id,
	}
}

func BuildTestNetQrCode(orderId, requester, payer string, from, to, value string) *tables.QrCode {
	return buildQrCode("Testnet", orderId, requester, payer, from, to, value)
}

func buildQrCode(chain, orderId, requester, payer string, from, to, value string) *tables.QrCode {
	exp := time.Now().Unix() + config.QrCodeExp
	data := &models.QrCodeData{
		Action: "transfer",
		Params: models.QrCodeParam{
			InvokeConfig: models.InvokeConfig{
				ContractHash: "",
				Functions: []models.Function{
					models.Function{
						Operation: "",
						Args: []models.Arg{
							models.Arg{
								Name:  "from",
								Value: from,
							},
							models.Arg{
								Name:  "to",
								Value: to,
							},
							models.Arg{
								Name:  "value",
								Value: value,
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
	return &tables.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		OrderId:    orderId,
		Requester:  requester,
		Signature:  "",
		Signer:     "",
		QrCodeData: string(databs),
		Callback:   config.DefConfig.QrCodeCallback,
		Exp:        exp,
		Chain:      chain,
		QrCodeDesc: "",
	}
}

func BuildApiBasicInfo(apiId int, icon, title, apiProvider, apiUrl, price, apiDesc string, specifications, popularity,
	delay, successRate, invokeFrequency int, createTime string) *tables.ApiBasicInfo {
	return &tables.ApiBasicInfo{
		ApiId:           apiId,
		Coin:            config.TOKEN_TYPE_ONG,
		ApiType:         config.Api,
		Icon:            icon,
		Title:           title,
		ApiProvider:     apiProvider,
		ApiUrl:          apiUrl,
		Price:           price,
		ApiDesc:         apiDesc,
		Specifications:  specifications,
		Popularity:      popularity,
		Delay:           delay,
		SuccessRate:     successRate,
		InvokeFrequency: invokeFrequency,
		CreateTime:      createTime,
	}
}
