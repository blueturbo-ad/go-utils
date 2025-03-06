package config_manage

import (
	"log"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
)

func TestHttpClientConfig(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	t.Run("test http client config file", func(t *testing.T) {
		var e = new(RedisConfigManager)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/http1_config.yaml"
		err = e.LoadConfig(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})

}
