package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	kvClient, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Slice to store all keys
	var keys []string

	log.Println("Getting all keys from Redis")
	// Get iterator
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()

	// Iterate over all keys
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}

	KVs := []*cloudflare.WorkersKVPair{}

	// Get values for all keys
	for _, key := range keys {
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				// Key does not exist
				continue
			}
			panic(err)
		}
		KVs = append(KVs, &cloudflare.WorkersKVPair{
			Key:   key,
			Value: val,
		})
	}

	log.Println("Writing keys to Cloudflare KV")
	// Store all keys in Cloudflare KV in bulk
	kvClient.WriteWorkersKVEntries(ctx, &cloudflare.ResourceContainer{
		Level:      "accounts",
		Identifier: os.Getenv("KV_USER_ID"),
		Type:       "account",
	}, cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: os.Getenv("KV_NAMESPACE_ID"),
		KVs:         KVs,
	},
	)
}
