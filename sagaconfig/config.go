package sagaconfig

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
)

var Version = ""

var (
	DEFAULT_LOG_LEVEL = log.InfoLog
	DEFAULT_REST_PORT = uint(8080)
)

type Config struct {
	NetWorkId           uint `json:"network_id"`
	OntSdk              *ontology_go_sdk.OntologySdk
	RestPort            uint      `json:"rest_port"`
	Version             string    `json:"version"`
	DbConfig            *DBConfig `json:"db_config"`
	OperatorPublicKey   string    `json:"operator_public_key"`
	ONTAuthScanProtocol string    `json:"ontauth_scan_protocol"`
	QrCodeCallback      string    `json:"qrcode_callback"`
	NASAAPIKey          string    `json:"nasa_api_key"`
	OntIdAccount        *ontology_go_sdk.Account
}

type DBConfig struct {
	ProjectDBUrl      string `json:"projectdb_url"`
	ProjectDBUser     string `json:"projectdb_user"`
	ProjectDBPassword string `json:"projectdb_password"`
	ProjectDBName     string `json:"projectdb_name"`
}

var DefDBConfigMap = map[int]*DBConfig{
	NETWORK_ID_SOLO_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1:3306",
		ProjectDBUser:     "root",
		ProjectDBPassword: "111111",
		ProjectDBName:     "saga",
	},
	NETWORK_ID_MAIN_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1:3306",
		ProjectDBUser:     "root",
		ProjectDBPassword: "111111",
		ProjectDBName:     "saga",
	},
	NETWORK_ID_TRAVIS_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1",
		ProjectDBUser:     "root",
		ProjectDBPassword: "",
		ProjectDBName:     "saga",
	},
}

var DefSagaConfig = &Config{
	RestPort:            DEFAULT_REST_PORT,
	Version:             "1.0.0",
	NetWorkId:           NETWORK_ID_POLARIS_NET,
	DbConfig:            DefDBConfigMap[NETWORK_ID_SOLO_NET],
	OperatorPublicKey:   "02b8fcf42deecc7cccb574ba145f2f627339fbd3ba2b63fda99af0a26a8d5a01da",
	ONTAuthScanProtocol: "http://192.168.1.175:8080/api/v1/onto/getQrCodeDataByQrCodeId",
	QrCodeCallback:      "http://192.168.1.175:8080/api/v1/onto/sendTx",
	NASAAPIKey:          NASA_API_KEY,
}
