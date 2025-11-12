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
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()

	k8sclient.GetSingleton().SetUp()
	t.Run("test EnvironmentConfig", func(t *testing.T) {
		logger := GetSingleton()
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		var p = fmt.Sprintf("%s//..//config//loggerex_conf.yaml", dir)
		fmt.Println("Current working directory:", p)

		err = logger.UpdateFromFile(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

		logger.Info("handle_logger_1", "{\"adx\":\"domob\",\"bid\":\"d43996a5-0198-1000-d01b-7d0717700030\",\"bid_time\":\"2025-08-23 08:00:02\",\"campaign_id\":622,\"creative_id\":8283,\"loss_content\":\"Creative Filtered - Disapproved by Exchange\",\"lossrsn\":\"202\",\"req_id\":\"63cfd00f50dcd-5f3b8bb1-11155157\"}")
		logger.Info("handle_logger_3", "{\"adx\":\"domob\",\"bid\":\"d43996a5-0198-1000-d01b-7d0717700030\",\"bid_time\":\"2025-08-23 08:00:02\",\"campaign_id\":622,\"creative_id\":8283,\"loss_content\":\"Creative Filtered - Disapproved by Exchange\",\"lossrsn\":\"202\",\"req_id\":\"63cfd00f50dcd-5f3b8bb1-11155157\"}")

		// logger.Info("handle_logger_1", "2222")
		// logger.Warn("handle_logger_1", "3333")
		// logger.Error("handle_logger_1", "4444")

		// logger.Info("system_logger", "2222")
		// logger.Warn("system_logger", "3333")
		// logger.Error("system_logger", "4444")

	})

}
