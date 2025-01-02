package config_manage

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestLoadEtcdConfig(t *testing.T) {
	// Setup a sample Config instance
	t.Run("test EnvironmentConfig", func(t *testing.T) {
		var e = new(EtcdifyConfig)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		var p = fmt.Sprintf("%s\\..\\config\\etcdify_conf.yaml", dir)
		err = e.LoadConfig(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})
}
