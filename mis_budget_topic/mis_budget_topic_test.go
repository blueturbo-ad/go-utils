package misbudgettopic

import (
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestMisBudGetTopic(t *testing.T) {

	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestMisBudGetTopic", func(t *testing.T) {
		a := GetSingleton()
		err := a.UpdateLoadK8sConfigMap("mis-to-budget-kafka-topic", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		conf := a.GetConfig("mis_budget")
		if a.Config == nil {
			t.Errorf("os.Getwd() = %v; want nil", a.Config)
		}
		if conf.Name != "mis_budget" {
			t.Errorf("os.Getwd() = %v; want nil", conf.Topic)
		}

	})

}
