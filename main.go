package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	redisclient "github.com/blueturbo-ad/go-utils/redis_client"
)

func main() {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	var e = redisclient.GetSingleton()
	err := e.UpdateLoadK8sConfigMap("redis-conf", "Dev")
	if err != nil {
		fmt.Println(err.Error())
	}
	client := e.GetReadClient("event_redis")
	if client == nil {
		fmt.Println(err.Error())
	}
	ctx := context.Background()
	res, err := client.Get(ctx, "test").Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}
