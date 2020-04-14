package core

import (
	"fmt"
	"github.com/candybox-sig/log"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	common2 "github.com/ontio/sagapi/common"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/dao"
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
	orderId, err := dao.DefDB.QueryOrderIdByQrCodeId(param.ExtraData.Id)
	if err != nil {
		return err
	}
	hash, err := config.DefConfig.OntSdk.SendTransaction(mutTx)
	if err != nil {
		return err
	}
	err = verifyTx(hash.ToHexString())
	if err != nil {
		return err
	}
	return dao.DefDB.UpdateTxInfoByOrderId(orderId, hash.ToHexString(), config.Completed)
}

func verifyTx(txHash string) error {
	retry := 0
	for {
		if retry > config.VERIFY_TX_RETRY {
			return fmt.Errorf("verify tx failed, txHash: %s", txHash)
		}
		event, err := config.DefConfig.OntSdk.GetSmartContractEvent(txHash)
		if err != nil {
			log.Errorf("[verifyTx] GetSmartContractEvent failed: %s", err)
			retry += 1
			continue
		}
		if event != nil && event.State == 1 {
			return nil
		}
		return fmt.Errorf("verify tx failed, txHash: %s", txHash)
	}
}
