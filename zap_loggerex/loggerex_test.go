package zap_loggerex

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

		logger.Info("handle_logger_1", "2222")
		logger.Warn("handle_logger_1", "3333")
		logger.Error("handle_logger_1", "4444")

	})

	t.Run("testBaseLoggerConfig1", func(t *testing.T) {
		// var e = new(ZapLoggerConfig)
		// e.LoadK8sConfigMap("dsp-logger", "Pro")
		logger := GetSingleton()
		logger.UpdateLoadK8sConfigMap("event-logger", "Pro")
		logger.GetCurConfig()
		for {
			time.Sleep(1 * time.Second)
			logger.Info("system_logger", "{\"level\":\"info\",\"ts\":\"2025-01-23T16:32:57.015+0800\",\"caller_line\":\"adapter/algorix_adapter.go:57\",\"msg\":\"62c5b75028035-167b945f-0006\\u00011737621177\\u00011737621177\\u0001/api/v1/dsp/Algorix\\u0001{\\\"id\\\":\\\"f645d6b545754dfeba5cc1640f7bb512\\\",\\\"imp\\\":[{\\\"id\\\":\\\"1\\\",\\\"banner\\\":{\\\"format\\\":[{\\\"w\\\":320,\\\"h\\\":50}],\\\"w\\\":320,\\\"h\\\":50,\\\"mimes\\\":[\\\"image/jpg\\\",\\\"image/jpeg\\\",\\\"image/gif\\\",\\\"image/png\\\"]},\\\"tagid\\\":\\\"com.campmobile.snow-BANNER\\\",\\\"bidfloor\\\":1.0693,\\\"bidfloorcur\\\":\\\"USD\\\",\\\"secure\\\":1}],\\\"app\\\":{\\\"id\\\":\\\"a99b37d8f8f9e668fe4d6c386a325587\\\",\\\"name\\\":\\\"SNOW - AI Profile\\\",\\\"bundle\\\":\\\"com.campmobile.snow\\\",\\\"storeurl\\\":\\\"https://play.google.com/store/apps/details?id=com.campmobile.snow\\\",\\\"publisher\\\":{\\\"id\\\":\\\"36870\\\"}},\\\"device\\\":{\\\"ua\\\":\\\"Mozilla/5.0 (Linux; Android 14; SM-A528B Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/131.0.6778.260 Mobile Safari/537.36\\\",\\\"geo\\\":{\\\"lat\\\":13.7063,\\\"lon\\\":100.4597,\\\"type\\\":2,\\\"country\\\":\\\"THA\\\",\\\"region\\\":\\\"10\\\",\\\"city\\\":\\\"Bangkok\\\",\\\"zip\\\":\\\"10110\\\"},\\\"ip\\\":\\\"49.229.170.154\\\",\\\"devicetype\\\":4,\\\"make\\\":\\\"samsung\\\",\\\"model\\\":\\\"SM-A528B\\\",\\\"os\\\":\\\"android\\\",\\\"osv\\\":\\\"14\\\",\\\"h\\\":853,\\\"w\\\":384,\\\"js\\\":1,\\\"language\\\":\\\"TH\\\",\\\"carrier\\\":\\\"00000\\\",\\\"connectiontype\\\":6,\\\"ifa\\\":\\\"0870bef6-8d57-42ef-9117-dc28aeacdd1b\\\"},\\\"user\\\":{\\\"ext\\\":{\\\"consent\\\":\\\"\\\"}},\\\"at\\\":1,\\\"tmax\\\":470,\\\"cur\\\":[\\\"USD\\\"],\\\"source\\\":{},\\\"regs\\\":{},\\\"exchange\\\":\\\"Algorix\\\"}\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001<nil>\\u0001\\u00011.07\\u00010\\u0001com.campmobile.snow\\u0001<nil>\\u0001<nil>\\u000110.153.0.151\\u0001\\u0001\\u00011\\u0001\\u0001\\u0001\\u0001\"}")
		}

	})

}
