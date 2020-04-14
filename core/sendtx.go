package core

import (
	"fmt"
	"github.com/candybox-sig/log"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/sagapi/config"
)

func SendTX(txHex string) error {
	txHexBs, err := common.HexToBytes(txHex)
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
	hash, err := config.DefConfig.OntSdk.SendTransaction(mutTx)
	if err != nil {
		return err
	}
	return verifyTx(hash.ToHexString())
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
