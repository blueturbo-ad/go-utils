package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// func main() {
// 	os.Setenv("POD_NAME", "test")
// 	os.Setenv("POD_NAMESPACE", "dsp-ns")
// 	environment.Init()
// 	workPath, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	cs := workPath + "/config/cloud_storage.yaml"
// 	err = gcpcloudstorage.GetSingleton().UpdateFromFile(cs, "Dev")
// 	if err != nil {
// 		fmt.Println("gcp cloud storage", err.Error())
// 	}
// 	ac := workPath + "/config/gcp_acc_token_conf.yaml"
// 	err = gcpcloudtool.GetSingleton().UpdateFromFile(ac, "Dev")
// 	if err != nil {
// 		fmt.Println("gcp access token", err.Error())
// 	}
// 	var e = redisclient.GetSingleton()

// 	p := workPath + "/config/redis_conf.yaml"
// 	err = e.UpdateFromFile(p, "Dev")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	client := e.GetReadClient("event_redis")
// 	if client == nil {
// 		fmt.Println("redis get client", err.Error())
// 	}
// 	ctx := context.Background()
// 	pong, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		log.Fatalf("Error connecting to Redis: %v", err)
// 	}
// 	fmt.Println(pong)
// }

func main() {

	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"10.0.0.215:6390",
			"10.0.0.215:6391",
			"10.0.0.215:6392",
			"10.0.0.215:6393",
			"10.0.0.215:6394",
			"10.0.0.215:6395",
		},
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		PoolSize:        10,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 30 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
	ctx := context.Background()
	pong, err := cli.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v", err)
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
	res, err := cli.Set(context.Background(), "key", "value", time.Duration(0)).Result()
	if err != nil {
		panic(err)
	}
	println(res)
	val, err := cli.Get(context.Background(), "key").Result()
	if err != nil {
		panic(err)
	}
	println(val)
}
