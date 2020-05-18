package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/types"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/core/http"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"strings"
	"time"
)

func SendTX(param *common2.SendTxParam) error {
	arr := strings.Split(param.Signer, ":")
	if len(arr) < 3 {
		return fmt.Errorf("error ontid: %s", param.Signer)
	}
	OwnerAddress := arr[2]

	txHexBs, err := common.HexToBytes(param.SignedTx)
	if err != nil {
		return err
	}
	tx := types.Transaction{}
	err = tx.Deserialization(common.NewZeroCopySource(txHexBs))
	if err != nil {
		return err
	}
	txHash := tx.Hash()
	log.Debugf("txHash: %s, signedTx: %s", txHash.ToHexString(), param.SignedTx)
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return err
	}
	order, err := dao.DefSagaApiDB.QueryOrderByQrCodeId(nil, param.ExtraData.Id)
	if err != nil {
		log.Debugf("SendTx.N.0: %s", err)
		return err
	}

	txdb, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("SendTx.N.1: %s", err)
		return errl
	}

	defer func() {
		if errl != nil {
			txdb.Rollback()
		}
	}()

	hash, err := sagaconfig.DefSagaConfig.OntSdk.SendTransaction(mutTx)
	if err != nil {
		return err
	}
	err = verifyTx(hash.ToHexString())
	if err != nil {
		log.Errorf("verifyTx failed: %s", err)
		err2 := dao.DefSagaApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, "", "", sagaconfig.Failed)
		if err2 != nil {
			log.Debugf("SendTx.N.10 %v", err)
			errl = err2
			return err2
		}
		return err
	} else {
		err = dao.DefSagaApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, "", hash.ToHexString(), sagaconfig.TxHandled)
		if err != nil {
			log.Debugf("SendTx.N.11 %v", err)
			errl = err
			return err
		}
	}

	result := ""
	if order.OrderType == sagaconfig.Api {
		err = generateApiKey(txdb, order.OrderId, param.Signer)
		if err != nil {
			log.Debugf("SendTx.N.2: %s", err)
			errl = err
			return err
		}
	} else if order.OrderType == sagaconfig.ApiProcess {
		// send request to server. check the result. and update status.
		// may check other request like coin. now default is WetherForcastServiceRequest
		api, err := dao.DefSagaApiDB.QueryApiBasicInfoByApiId(txdb, order.ApiId, tables.API_STATE_BUILTIN)
		if err != nil {
			errl = err
			log.Debugf("SendTx.N.3: %s", err)
			return err
		}

		if order.OrderKind == tables.ORDER_KIND_DATA_PROCESS_WETHER {
			paramWether := &common2.WetherForcastRequest{}
			err = json.Unmarshal([]byte(order.Request), paramWether)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.4: %s", err)
				return err
			}
			env, err := dao.DefSagaApiDB.QueryEnvById(txdb, paramWether.EnvId)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.5: %s", err)
				return err
			}
			alg, err := dao.DefSagaApiDB.QueryAlgorithmById(nil, paramWether.AlgorithmId)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.6: %s", err)
				return err
			}
			tm := time.Unix(paramWether.TargetDate, 0)
			mm, err := time.ParseDuration("-240h")
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.7 %v", err)
				return err
			}
			targetTime := tm.Add(mm)

			r := common2.WetherForcastServiceRequest{
				DataUrl:       api.ApiProvider,
				Header:        make(map[string]interface{}),
				Param:         make(map[string]interface{}),
				RequestMethod: api.RequestType,
				AlgorithmName: alg.AlgName,
				Dtoken: common2.Dtoken{
					ResourceId: api.ResourceId,
					Account:    OwnerAddress,
					TokenTemplate: common2.TokenTemplate{
						DataIds:   "",
						TokenHash: api.TokenHash,
					},
					Number: 1,
				},
			}

			r.Header["Authorization"] = "4cdb5582-90d8-11ea-af71-0242ac130002-4cdb5640-90d8-11ea-af71-0242ac130002"
			r.Param["params"] = "airTemperature"
			r.Param["lat"] = paramWether.Location.Lat
			r.Param["lng"] = paramWether.Location.Lng
			r.Param["start"] = targetTime.Format("2006-01-02") // should be string.
			log.Debugf("SendTx.Y.0 %v", r)

			data, err := json.Marshal(r)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.8 %v", err)
				return err
			}
			res, err := http.NewClient().Post(env.ServiceUrl, data)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.9 %v", err)
				return err
			}

			type DataProcessServiceRespone struct {
				Action  string      `json:"action"`
				Error   uint32      `json:"error"`
				Desc    string      `json:"desc"`
				Result  interface{} `json:"result"`
				Version string      `json:"version"`
			}
			var t DataProcessServiceRespone
			err = json.Unmarshal(res, &t)
			if err != nil {
				errl = err
				log.Debugf("SendTx.N.9.0 %s", err)
				return err
			}

			if t.Error != 0 {
				errl = errors.New(t.Desc)
				log.Debugf("SendTx.N.9.1 %s", errl)
				return err
			}

			log.Debugf("%s: %s", alg.AlgName, string(res))
			result = string(res)
		} else {
			errl = errors.New("wrong data process type")
			return errl
		}
	}

	err = dao.DefSagaApiDB.UpdateTxInfoByOrderId(txdb, order.OrderId, result, hash.ToHexString(), sagaconfig.Completed)
	if err != nil {
		log.Debugf("SendTx.N.11 %v", err)
		errl = err
		return err
	}

	err = txdb.Commit()
	if err != nil {
		errl = err
		log.Debugf("SendTx.N.12 %v", err)
		return err
	}
	return nil
}

func generateApiKey(tx *sqlx.Tx, orderId, ontId string) error {
	order, err := dao.DefSagaApiDB.QueryOrderByOrderId(tx, orderId)
	if err != nil {
		return err
	}

	spec, err := dao.DefSagaApiDB.QuerySpecificationsById(tx, order.SpecificationsId)
	if err != nil {
		return err
	}

	id := common2.GenerateUUId(common2.UUID_TYPE_API_KEY)
	apiKey := &tables.APIKey{
		OrderId:      orderId,
		ApiKey:       id,
		ApiId:        order.ApiId,
		RequestLimit: spec.Amount,
		UsedNum:      0,
		OntId:        ontId,
	}
	return dao.DefSagaApiDB.InsertApiKey(tx, apiKey)
}

func verifyTx(txHash string) error {
	retry := 0
	for {
		if retry > sagaconfig.VERIFY_TX_RETRY {
			return fmt.Errorf("verify tx failed, txHash: %s", txHash)
		}
		event, err := sagaconfig.DefSagaConfig.OntSdk.GetSmartContractEvent(txHash)
		if err != nil && strings.Contains(err.Error(), "duplicated transaction detected") {
			return err
		}
		if err != nil || event == nil {
			log.Errorf("[verifyTx] GetSmartContractEvent failed: %s", err)
			sleepTime := sagaconfig.VERIFY_TX_RETRY - retry
			time.Sleep(time.Duration(sleepTime) * time.Second)
			retry += 1
			log.Infof("[verifyTx] txHash: %s, retry:%d, err: %s", txHash, retry, err)
			continue
		}
		if event != nil && event.State == 1 {
			log.Infof("txHash:%s, event.State:%d", txHash, event.State)
			return nil
		}
		return fmt.Errorf("verify tx failed, txHash: %s", txHash)
	}
}
