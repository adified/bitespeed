package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_Url string `mapstructure:"db_url"`
}

func Loadconfig() (Config, error) {
	var C Config

	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/")   // path to look for the config file in
	viper.AddConfigPath(".")   // optionally look for config in the working directory

	// err := viper.ReadInConfig() // Find and read the config file
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return C, err
	}
	err = viper.Unmarshal(&C)

	if err != nil { // Handle errors reading the config file
		fmt.Printf("log.Fatal: ")
	}
	fmt.Println()

	return C, nil
}
