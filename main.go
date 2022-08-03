package main // import "github.com/italia/developers-italia-backend"

import (
	"fmt"

	"github.com/italia/developers-italia-backend/cmd"
	"github.com/italia/developers-italia-backend/crawler"

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
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error reding config file: %s", err))
	}

	// Register client APIs.
	crawler.RegisterClientAPIs()

	cmd.Execute()
}
