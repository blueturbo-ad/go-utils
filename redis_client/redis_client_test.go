package redisclient

import (
	"context"
	"fmt"
	"os"
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
}
