package dspbaseconfig

import (
	"fmt"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestDspBaseConfig(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestDspBaseConfig", func(t *testing.T) {
		workPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}

		cs := workPath + "/config/billing-model-config.yaml"
		fmt.Println(cs)
		GetSingleton().RegistHookFunc(func(config string) error {
			fmt.Println("hook func", config)
			return nil
		})
		err = GetSingleton().LoadFileConfig(cs, "Dev", "")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})
}
