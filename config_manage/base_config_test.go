package config_manage

import (
	"testing"

	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	"gopkg.in/yaml.v3"
)

type Configs struct {
	CurUsed string      ` yaml:"curused"`
	Dev     interface{} `yaml:"Dev"`
	Pro     interface{} `yaml:"Pro"`
	Pre     interface{} `yaml:"Pre"`
	Test    interface{} `yaml:"Test"`
}

func TestBaseConfig(t *testing.T) {

	t.Run("TestBaseFeishuConfig", func(t *testing.T) {
		c := new(ManagerConfig)
		k8sclient.GetSingleton().SetUp()
		c.LoadK8sConfigMap("dsp-ns", "feishu", "Dev")
	})
	t.Run("TestBaseManagerConfig", func(t *testing.T) {
		str := `
curused: "Dev"
Dev:
  url: "https://open.feishu.cn/open-apis/bot/v2/hook/7033a584-ce85-4eea-ae9d-fe1f5589ab59"
Pro:
  url: "https://open.feishu.cn/open-apis/bot/v2/hook/7033a584-ce85-4eea-ae9d-fe1f5589ab59"
Pre:
  url: "https://open.feishu.cn/open-apis/bot/v2/hook/7033a584-ce85-4eea-ae9d-fe1f5589ab59"
`

		var config = Configs{}
		err := yaml.Unmarshal([]byte(str), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal YAML: %v", err)
		}
	})
}
