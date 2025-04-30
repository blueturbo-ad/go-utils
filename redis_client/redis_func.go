package redisclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// 获取redis hSet中的ckey
func GetRedisHashString(cli *redis.ClusterClient, key string, ckey string) (string, error) {
	if cli == nil {
		return "", fmt.Errorf("redis ClusterClient: redis client is nil")
	}
	// 获取 hget key OfferMsg
	allmap, err := cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("GetRedisMessage: redis hgetall key:%s error %s", key, err.Error())
	}
	// 获取 conv_cap 固定值
	value, ok := allmap[ckey]
	if !ok {
		return "", fmt.Errorf("GetRedisMessage: conv_cap not found in redis ckey:%s", ckey)
	}

	return value, nil
}

// 获取redis hSet中的ckey
func GetRedisHashUint64(cli *redis.ClusterClient, key string, ckey string) (uint64, error) {
	if cli == nil {
		return 0, fmt.Errorf("redis ClusterClient: redis client is nil")
	}
	// 获取 hget key OfferMsg
	allmap, err := cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("GetRedisMessage: redis hgetall key:%s error %s", key, err.Error())
	}
	// 获取 conv_cap 固定值
	value, ok := allmap[ckey]
	if !ok {
		return 0, fmt.Errorf("GetRedisMessage:  not found in redis ckey:%s", ckey)
	}
	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("GetRedisMessage: strconv.ParseUint parse error %s", err.Error())
	}

	return num, nil
}

func GetRedisHashFloat(cli *redis.ClusterClient, key string, ckey string) (float64, error) {
	if cli == nil {
		return 0, fmt.Errorf("redis ClusterClient: redis client is nil")
	}
	// 获取 hget key OfferMsg
	allmap, err := cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("GetRedisMessage: redis hgetall key:%s error %s", key, err.Error())
	}
	// 获取 conv_cap 固定值
	value, ok := allmap[ckey]
	if !ok {
		return 0, fmt.Errorf("GetRedisMessage:  not found in redis ckey:%s", ckey)
	}
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("GetRedisMessage: strconv.ParseUint parse error %s", err.Error())
	}

	return num, nil
}
