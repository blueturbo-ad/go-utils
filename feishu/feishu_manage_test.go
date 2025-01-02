package feishu

import (
	"log"
	"os"
	"testing"
)

func TestFeishuManage(t *testing.T) {
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
}
