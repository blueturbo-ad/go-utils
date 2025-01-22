package zap_loggerex

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

// func BenchmarkUpdateFromFile(b *testing.B) {
// 	logger := GetSingleton()
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var p = fmt.Sprintf("%s/loggerex_conf.yaml", dir)
// 	fmt.Println("Current working directory:", p)
// 	logger.UpdateFromFile(p, "Dev")
// 	for i := 0; i < b.N; i++ {
// 		logger.Info("handle_logger_1", "2222")
// 		logger.Warn("handle_logger_1", "3333")
// 	}

// }

func TestGetLogger(t *testing.T) {
	// Setup a sample Config instance
	environment.Init()
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	k8sclient.GetSingleton().SetUp()
	t.Run("test EnvironmentConfig", func(t *testing.T) {
		logger := GetSingleton()
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		var p = fmt.Sprintf("%s\\..\\config\\loggerex_conf.yaml", dir)
		fmt.Println("Current working directory:", p)

		err = logger.UpdateFromFile(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

		logger.Info("handle_logger_1", "2222")
		logger.Warn("handle_logger_1", "3333")
		logger.Error("handle_logger_1", "4444")

		logger.Debug("handle_logger_2", "1111")
		logger.Info("handle_logger_2", "2222")
		logger.Warn("handle_logger_2", "3333")
		logger.Error("handle_logger_2", "4444")

		logger.Debug("handle_logger_3", "1111")
		logger.Info("handle_logger_3", "2222")
		logger.Warn("handle_logger_3", "3333")
		logger.Error("handle_logger_3", "4444")

	})

	t.Run("testBaseLoggerConfig1", func(t *testing.T) {
		// var e = new(ZapLoggerConfig)
		// e.LoadK8sConfigMap("dsp-logger", "Pro")
		logger := GetSingleton()
		logger.UpdateLoadK8sConfigMap("dsp-logger", "Pro")
		logger.Debug("handle_logger_1", "1111")
	})

}
