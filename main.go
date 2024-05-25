package main

import (
	"errors"
	"fmt"

	"github.com/italia/publiccode-crawler/v4/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// Read configurations.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetDefault("DATADIR", "./data")
	viper.SetDefault("ACTIVITY_DAYS", 60)
	viper.SetDefault("API_BASEURL", "https://api.developers.italia.it/v1/")
	viper.SetDefault("MAIN_PUBLISHER_ID", "")
	viper.SetDefault("GITHUB_TOKEN", "")

	if err := viper.ReadInConfig(); err != nil {
		var notFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundError) {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	}

	cmd.Execute()
}
