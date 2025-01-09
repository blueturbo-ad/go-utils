package config_manage

import (
	"testing"

	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestBaseConfig(t *testing.T) {

	t.Run("TestBaseFeishuConfig", func(t *testing.T) {
		c := new(ManagerConfig)
		k8sclient.GetSingleton().SetUp()
		c.LoadK8sConfigMap("dsp-ns", "feishu", "Dev")
	})

}
