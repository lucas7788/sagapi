package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

func SetOntologyConfig(ctx *cli.Context) error {
	cf := ctx.String(GetFlagName(ConfigfileFlag))
	if _, err := os.Stat(cf); os.IsNotExist(err) {
		// if there's no config file, use default config
		updateConfigByCmd(ctx)
		return nil
	}
	file, err := os.Open(cf)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	cfg := &sagaconfig.Config{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return err
	}
	*sagaconfig.DefSagaConfig = *cfg

	sdk := ontology_go_sdk.NewOntologySdk()
	switch sagaconfig.DefSagaConfig.NetWorkId {
	case sagaconfig.NETWORK_ID_MAIN_NET:
		sdk.NewRpcClient().SetAddress(sagaconfig.ONT_MAIN_NET)
		sagaconfig.DefSagaConfig.NetType = sagaconfig.MainNet
	case sagaconfig.NETWORK_ID_POLARIS_NET:
		sdk.NewRpcClient().SetAddress(sagaconfig.ONT_TEST_NET)
		sagaconfig.DefSagaConfig.NetType = sagaconfig.TestNet
	case sagaconfig.NETWORK_ID_SOLO_NET:
		sdk.NewRpcClient().SetAddress(sagaconfig.ONT_SOLO_NET)
	default:
		return fmt.Errorf("error network id %d", sagaconfig.DefSagaConfig.NetWorkId)
	}

	sagaconfig.DefSagaConfig.OntSdk = sdk
	return nil
}

func updateConfigByCmd(ctx *cli.Context) error {
	port := ctx.Uint(GetFlagName(RestPortFlag))
	if port != 0 {
		sagaconfig.DefSagaConfig.RestPort = port
	}
	networkId := ctx.Uint(GetFlagName(NetworkIdFlag))
	if networkId > 3 {
		return fmt.Errorf("networkid should be between 1 and 3, curr: %d", networkId)
	}
	sagaconfig.DefSagaConfig.NetWorkId = networkId
	//rpc := sagaconfig.ONT_MAIN_NET
	//if networkId == sagaconfig.NETWORK_ID_POLARIS_NET {
	//	rpc = sagaconfig.ONT_TEST_NET
	//} else if networkId == sagaconfig.NETWORK_ID_SOLO_NET {
	//	rpc = sagaconfig.ONT_SOLO_NET
	//}
	sdk := ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(sagaconfig.ONT_TEST_NET)
	sagaconfig.DefSagaConfig.OntSdk = sdk
	return nil
}

func PrintErrorMsg(format string, a ...interface{}) {
	format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format) //Print error msg with red color
	fmt.Printf(format, a...)
}
