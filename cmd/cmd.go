package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/sagapi/config"
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
	cfg := &config.Config{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return err
	}
	*config.DefSagaConfig = *cfg
	updateConfigByCmd(ctx)
	return nil
}

func updateConfigByCmd(ctx *cli.Context) error {
	port := ctx.Uint(GetFlagName(RestPortFlag))
	if port != 0 {
		config.DefSagaConfig.RestPort = port
	}
	networkId := ctx.Uint(GetFlagName(NetworkIdFlag))
	if networkId > 3 {
		return fmt.Errorf("networkid should be between 1 and 3, curr: %d", networkId)
	}
	config.DefSagaConfig.NetWorkId = networkId
	rpc := config.ONT_MAIN_NET
	if networkId == config.NETWORK_ID_POLARIS_NET {
		rpc = config.ONT_TEST_NET
	} else if networkId == config.NETWORK_ID_SOLO_NET {
		rpc = config.ONT_SOLO_NET
	}
	fmt.Println("rpc:", rpc)
	sdk := ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(config.ONT_TEST_NET)
	config.DefSagaConfig.OntSdk = sdk
	return nil
}

func PrintErrorMsg(format string, a ...interface{}) {
	format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format) //Print error msg with red color
	fmt.Printf(format, a...)
}
