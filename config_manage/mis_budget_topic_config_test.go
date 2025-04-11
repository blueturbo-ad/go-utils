package config_manage

import (
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestMisBudgetTopicConfig(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestMisBudgetTopicConfig", func(t *testing.T) {
		var e = new(GcpSvcAccountTokenConfig)
		err := e.LoadK8sConfigMap("mis-to-budget-kafka-topic", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})
}
