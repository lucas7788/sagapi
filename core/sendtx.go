package core

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/types"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/dao"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"strings"
	"time"
)

func SendTX(param *common2.SendTxParam) error {
	txHexBs, err := common.HexToBytes(param.SignedTx)
	if err != nil {
		return err
	}
	tx := types.Transaction{}
	err = tx.Deserialization(common.NewZeroCopySource(txHexBs))
	if err != nil {
		return err
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return err
	}
	hash, err := sagaconfig.DefSagaConfig.OntSdk.SendTransaction(mutTx)
	if err != nil {
		return err
	}
	orderId, err := dao.DefSagaApiDB.QrCodeDB.QueryOrderIdByQrCodeId(param.ExtraData.Id)
	if err != nil {
		return err
	}
	err = verifyTx(hash.ToHexString())
	if err != nil {
		log.Errorf("verifyTx failed: %s", err)
		err2 := dao.DefSagaApiDB.OrderDB.UpdateTxInfoByOrderId(orderId, "", sagaconfig.Failed)
		if err2 != nil {
			return err2
		}
		return err
	}
	err = generateApiKey(orderId, param.Signer)
	if err != nil {
		return err
	}
	err = dao.DefSagaApiDB.OrderDB.UpdateTxInfoByOrderId(orderId, hash.ToHexString(), sagaconfig.Completed)
	if err != nil {
		return err
	}
	return nil
}

func generateApiKey(orderId, ontId string) error {
	order, err := dao.DefSagaApiDB.OrderDB.QueryOrderByOrderId(orderId)
	if err != nil {
		return err
	}
	spec, err := dao.DefSagaApiDB.ApiDB.QuerySpecificationsBySpecificationsId(order.SpecificationsId)
	if err != nil {
		return err
	}
	id := common2.GenerateUUId()
	apiKey := &tables.APIKey{
		OrderId:      orderId,
		ApiKey:       id,
		ApiId:        order.ApiId,
		RequestLimit: spec.Amount,
		UsedNum:      0,
		OntId:        ontId,
	}
	return dao.DefSagaApiDB.ApiDB.InsertApiKey(apiKey)
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
