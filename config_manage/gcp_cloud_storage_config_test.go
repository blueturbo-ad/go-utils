package config_manage

import (
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestGcpCloudStorageTok(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestGcpSvcAccTok", func(t *testing.T) {
		var e = new(GcpCloudStorageTokenConfig)
		err := e.LoadK8sConfigMap("gcp-cloud-storage-config", "Dev")
		if err != nil {
			t.Errorf("os.Getwd() = %v; want nil", err)
		}
	})
}
