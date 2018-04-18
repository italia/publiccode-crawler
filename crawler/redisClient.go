package crawler

import (
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
)

func redisClientFactory(URL string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     URL, // docker redis ip "redis:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
	})

	// Check if connection is available.
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Errorf("error on redisClientFactory: %s:", err)
		return nil, err
	}

	return redisClient, nil

}
