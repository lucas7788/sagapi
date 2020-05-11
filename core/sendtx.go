package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
	txHash := tx.Hash()
	log.Debugf("txHash: %s, signedTx: %s", txHash.ToHexString(), param.SignedTx)
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return err
	}
	hash, err := sagaconfig.DefSagaConfig.OntSdk.SendTransaction(mutTx)
	if err != nil {
		return err
	}

	txdb, errl := dao.DefSagaApiDB.DB.Beginx()
	if errl != nil {
		return errl
	}

	defer func() {
		if errl != nil {
			txdb.Rollback()
		}
	}()

	orderId, err := dao.DefSagaApiDB.QueryOrderIdByQrCodeId(nil, param.ExtraData.Id)
	if err != nil {
		return err
	}
	err = verifyTx(hash.ToHexString())
	if err != nil {
		log.Errorf("verifyTx failed: %s", err)
		err2 := dao.DefSagaApiDB.UpdateTxInfoByOrderId(txdb, orderId, "", sagaconfig.Failed)
		if err2 != nil {
			errl = err2
			return err2
		}
		return err
	}
	err = generateApiKey(txdb, orderId, param.Signer)
	if err != nil {
		errl = err
		return err
	}
	err = dao.DefSagaApiDB.UpdateTxInfoByOrderId(txdb, orderId, hash.ToHexString(), sagaconfig.Completed)
	if err != nil {
		errl = err
		return err
	}
	err = txdb.Commit()
	errl = err
	return err
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

	id := common2.GenerateUUId(common2.UUID_TYPE_RAW)
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
