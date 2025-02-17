package config_manage

import (
	"log"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestRredis(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("text redis config file", func(t *testing.T) {
		var e = new(RedisConfigManager)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/redis_conf.yaml"
		err = e.LoadConfig(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})

	t.Run("test redis configmap", func(t *testing.T) {
		var e = new(RedisConfigManager)
		err := e.LoadK8sConfigMap("redis-conf", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})

}
