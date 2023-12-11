package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	var app ProxyApp
	viper.SetConfigName("properties")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&app)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	app.Init()
}