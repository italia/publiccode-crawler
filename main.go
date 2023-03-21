package main

import (
	"fmt"

	"github.com/italia/publiccode-crawler/v3/cmd"
	"github.com/spf13/viper"
)

func main() {
	// Read configurations.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetDefault("DATADIR", "./data")
	viper.SetDefault("ACTIVITY_DAYS", 60)
	viper.SetDefault("API_BASEURL", "https://api.developers.italia.it/v1/")
	viper.SetDefault("MAIN_PUBLISHER_ID", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, fileNotFound := err.(viper.ConfigFileNotFoundError); !fileNotFound {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	}

	cmd.Execute()
}
