package common

import (
	"encoding/hex"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models"
	"github.com/ontio/sagapi/models/tables"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

func GenerateOrderId() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func BuildTestNetQrCode(requester, payer string, from, to, value string) *tables.QrCode {
	return buildQrCode("Testnet", requester, payer, from, to, value)
}

func buildQrCode(chain, requester, payer string, from, to, value string) *tables.QrCode {
	now := time.Now().Nanosecond()
	exp := time.Now().Unix() + config.QrCodeExp
	data := &models.QrCodeData{
		Action: "transfer",
		Params: models.Param{
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
	id := hex.EncodeToString([]byte(strconv.Itoa(now)))
	return &tables.QrCode{
		Ver:       "1.0.0",
		Id:        id,
		Requester: requester,
		Signature: "",
		Signer:    "",
		Data:      data,
		Callback:  "http://127.0.0.1:8080/api/v1/sendtx",
		Exp:       exp,
		Chain:     chain,
		Desc:      "",
	}
}
