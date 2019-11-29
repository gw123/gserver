package gserver

import (
	"github.com/gw123/gserver/contracts"
	"github.com/spf13/viper"
)

func LoadServerConfig() contracts.ServerConfig {
	myConfig := viper.New()
	viper.AutomaticEnv()
	myConfig.SetConfigFile("config.server.yaml")
	//myConfig.SetConfigFile("/etc/gserver/config.server.yaml")
	if err := myConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := contracts.ServerConfig{}
	if err := myConfig.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	return serverConfig
}
