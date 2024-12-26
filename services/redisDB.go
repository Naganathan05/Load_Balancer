package services

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client

var ctx = context.Background()

func GetOptimalServer(zSetKey string) (string, error) {
	optimalServer, err := client.ZRangeWithScores(ctx, zSetKey, 0, 0).Result()
	if err != nil {
		return "", fmt.Errorf("Error in retrieving top server from ZSET : %w", err)
	}

	if len(optimalServer) == 0 {
		return "", fmt.Errorf("No Healthy server found : %s", zSetKey)
	}

	IP_Addr := optimalServer[0]

	return IP_Addr.Member.(string), nil
}

func UpdateActiveCount(zSetKey string, IP_Addr string, updateVal float64) (float64, error) {
	newVal, err := client.ZIncrBy(ctx, zSetKey, updateVal, IP_Addr).Result()
	if err != nil {
		return -1, err
	}
	return newVal, nil
}

func InitRedisClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	fmt.Println("Initialised Redis Client.")

	// opt, err:= redis.ParseURL("URL")
	// if err != nil {
	// 	panic(err)
	// }
	// client := redis.NewClient(opt)

	// data, err := GetOptimalServer("www.google.com")
	// if err != nil {
	// 	panic(err)
	// }

	// data, err := UpdateActiveCount("www.google.com", "192.168.2.2", -1)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Ip Address : %s", data)
}
