package config

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
	NetWorkId uint
	OntSdk    *ontology_go_sdk.OntologySdk
	RestPort  uint   `json:"rest_port"`
	Version   string `json:"version"`
	DbConfig  *DBConfig
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
}

var DefConfig = &Config{
	RestPort:  DEFAULT_REST_PORT,
	Version:   "1.0.0",
	NetWorkId: NETWORK_ID_SOLO_NET,
	DbConfig:  DefDBConfigMap[NETWORK_ID_SOLO_NET],
}
