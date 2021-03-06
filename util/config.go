package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config struct defines the config structure
type Config struct {
	Mongo MongoConfig `json:"mongo,omitempty"`
	Mysql MysqlConfig `json:"mysql,omitempty"`
	Host  string      `json:"host"`
}

// MongoConfig has config values for Mongo
type MongoConfig struct {
	Addr  string `json:"addr"`
	DB    string `json:"db"`
	Table string `json:"table"`
	Event string `json:"event"`
}

type MysqlConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DB     string `json:"db"`
	Table  string `json:"table"`
}

// NewConfig parses config file and return Config struct
func NewConfig() *Config {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("Read config file error.")
	}
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		panic(err)
	}

	return config
}
