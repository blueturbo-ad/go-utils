package redisclient

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestRedisClientTest(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestRedisClientTest", func(t *testing.T) {
		var e = GetSingleton()
		err := e.UpdateLoadK8sConfigMap("redis-conf", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})
	t.Run("TestRedisGetReadClient", func(t *testing.T) {
		var e = GetSingleton()
		err := e.UpdateLoadK8sConfigMap("redis-conf", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		client := e.GetReadClient("event_redis")
		if client == nil {
			t.Errorf("os.Getwd() = %v; want nil", client)
		}
		ctx := context.Background()
		res, err := client.Get(ctx, "test").Result()
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		fmt.Println(res)
	})
	t.Run("TestRedisGetWriteClient", func(t *testing.T) {
		workPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}
		// 去workpath 的上一级目录
		works := strings.Split(workPath, "/")
		workPath = strings.Join(works[:len(works)-1], "/")

		cs := workPath + "/config/redis_conf.yaml"
		err = GetSingleton().UpdateFromFile(cs, "Dev")
		if err != nil {
			fmt.Println(err.Error())
		}
		var e = GetSingleton()
		client := e.GetWriteClient("event_redis")
		if client == nil {
			t.Errorf("os.Getwd() = %v; want nil", client)
		}
		ctx := context.Background()
		res, err := client.Get(ctx, "test").Result()
		if err != nil {
			fmt.Println(err.Error())
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		fmt.Println(res)
	})
}
