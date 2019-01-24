package main // import "github.com/italia/developers-italia-backend/crawler"

import (
	"fmt"
	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/italia/developers-italia-backend/crawler/cmd"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// Read configurations.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reding config file: %s", err))
	}

	// Register client APIs.
	crawler.RegisterClientAPIs()

	cmd.Execute()
}
