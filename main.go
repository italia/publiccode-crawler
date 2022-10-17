package main // import "github.com/italia/developers-italia-backend"

import (
	"fmt"

	"github.com/italia/developers-italia-backend/cmd"

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

	viper.SetDefault("API_BASEURL", "https://api.developers.italia.it/v1/")

	if err := viper.ReadInConfig(); err != nil {
		if _, fileNotFound := err.(viper.ConfigFileNotFoundError); !fileNotFound {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	}

	cmd.Execute()
}
