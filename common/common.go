package common

import (
	"encoding/json"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/utils"
	"github.com/satori/go.uuid"
	"time"
)

func GenerateUUId() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func BuildQrCodeResponse(id string) *QrCodeResponse {
	return &QrCodeResponse{
		QrCode: QrCode{
			ONTAuthScanProtocol: config.DefSagaConfig.ONTAuthScanProtocol + "/" + id,
		},
		Id: id,
	}
}

func BuildTestNetQrCode(orderId, ontid, payer, from, to, value string) *tables.QrCode {
	return buildQrCode("Testnet", orderId, ontid, payer, from, to, value)
}

func buildQrCode(chain, orderId, ontid, payer, from, to, value string) *tables.QrCode {
	exp := time.Now().Unix() + config.QrCodeExp
	amt := utils.ToIntByPrecise(value, config.ONG_DECIMALS)
	data := &models.QrCodeData{
		Action: "signTransaction",
		Params: models.QrCodeParam{
			InvokeConfig: models.InvokeConfig{
				ContractHash: config.ONG_CONTRACT_ADDRESS,
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
	sig, err := config.DefSagaConfig.OntIdAccount.Sign(databs)
	if err != nil {

	}
	return &tables.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		OrderId:    orderId,
		Requester:  config.OntId,
		Signature:  common.ToHexString(sig),
		Signer:     ontid,
		QrCodeData: string(databs),
		Callback:   config.DefSagaConfig.QrCodeCallback,
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
