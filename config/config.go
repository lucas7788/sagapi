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
	NetWorkId         uint
	OntSdk            *ontology_go_sdk.OntologySdk
	RestPort          uint   `json:"rest_port"`
	Version           string `json:"version"`
	ProjectDBHost     string `json:"projectdb_host"`
	ProjectDBPort     string `json:"projectdb_port"`
	ProjectDBUrl      string `json:"projectdb_url"`
	ProjectDBUser     string `json:"projectdb_user"`
	ProjectDBPassword string `json:"projectdb_password"`
	ProjectDBName     string `json:"projectdb_name"`
}

var DefConfig = &Config{
	RestPort:  DEFAULT_REST_PORT,
	Version:   "1.0.0",
	NetWorkId: NETWORK_ID_MAIN_NET,
}
