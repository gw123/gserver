package config

import (
	"github.com/gw123/gserver/contracts"
	"github.com/spf13/viper"
)

func LoadClientConfig() contracts.ClientConfig {
	myConfig := viper.New()
	viper.AutomaticEnv()
	myConfig.SetConfigFile("/etc/gserver/config.client.yaml")
	if err := myConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := contracts.ClientConfig{}
	if err := myConfig.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	return serverConfig
}