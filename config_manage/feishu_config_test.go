package config_manage

import (
	"log"
	"os"
	"testing"
)

func TestFeishuConfig(t *testing.T) {
	t.Run("TestFeishuConfig", func(t *testing.T) {
		// Test code here
		var e = new(FeishuConfig)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/feishu_config.yaml"
		err = e.LoadConfig(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})
}
