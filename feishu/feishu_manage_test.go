package feishu

import (
	"log"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestFeishuManage(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestFeishuerrorManage", func(t *testing.T) {
		// Test code here
		var f = GetInstance()
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/feishu_config.yaml"
		err = f.UpdateFromFile(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		err = f.Send("error", "这是一条错误测试信息")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})
	t.Run("TestFeishuwarningManage", func(t *testing.T) {
		// Test code here
		var f = GetInstance()
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/feishu_config.yaml"
		err = f.UpdateFromFile(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		err = f.Send("warning", "这是一条warning测试信息")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})

	t.Run("TestFeishuFailManage", func(t *testing.T) {
		// Test code here
		var f = GetInstance()
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		p := dir + "/config/feishu_config.yaml"
		err = f.UpdateFromFile(p, "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
		err = f.Send("fail", "这是一条fail测试信息")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}

	})
	t.Run("TestFeishuConfigMap", func(t *testing.T) {
		feishu := GetInstance()
		feishu.UpdateLoadK8sConfigMap("feishu", "Pro", "")
	})
}
