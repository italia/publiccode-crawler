package cmd

import (
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(redisCmd)
}

var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Print the version number of the crawler",
	Long:  `All software has versions. This too.`,
	Run: func(cmd *cobra.Command, args []string) {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		// Output: PONG <nil> for test
		pong, err := redisClient.Ping().Result()
		fmt.Println(pong, err)

		keys, _ := redisClient.Keys("*").Result()
		for _, key := range keys {
			if redisClient.Get(key).Val() == "false" {
				log.Error("Found one false URL! start from here: " + key)
			}
		}
	},
}
