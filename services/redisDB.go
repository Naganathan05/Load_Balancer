package services

import (
	serverSchema "Load_Balancer_Server/models"
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
)

func getItem(client *redis.Client, domainName string) (serverSchema.ServerData, error) {
	ctx := context.Background()
	var val serverSchema.ServerData

	// data, err := client.Get(ctx, domainName).Result()
	// if err != nil {
	// 	return val, err
	// }

	// // Unmarshal the JSON string into the struct
	// err = json.Unmarshal([]byte(data), &val)
	// if err != nil {
	// 	return val, fmt.Errorf("error unmarshalling JSON: %w", err)
	// }
	// Use .Scan to decode directly

	err := client.Get(ctx, domainName).Scan(&val)
	if err != nil {
		return val, fmt.Errorf("error scanning Redis data: %w", err)
	}

	return val, nil
}

func RedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// opt, err:= redis.ParseURL("URL")
	// if err != nil {
	// 	panic(err)
	// }
	// client := redis.NewClient(opt)

	val, err := getItem(client, "www.google.com")
	if err != nil {
		panic(err)
	}

	fmt.Println(val.IpAddr)
}
