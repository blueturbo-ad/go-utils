package etcdify

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	logger "github.com/blueturbo-ad/go-utils/zap_loggerex"
)

func TestNewWatcher(t *testing.T) {
	// Setup a sample Config instance
	t.Run("test EnvironmentConfig", func(t *testing.T) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		var p = fmt.Sprintf("%s\\..\\config\\etcdify_conf.yaml", dir)

		etcd, err := NewWatcher(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
			return
		}
		loggerex := logger.GetSingleton()

		etcd.WatchKey("Dev", context.Background(), "loggerex", loggerex.UpdateFromEtcd)

	})
}
