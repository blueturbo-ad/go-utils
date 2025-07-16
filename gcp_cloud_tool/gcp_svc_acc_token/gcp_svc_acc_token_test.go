package gcpcloudtool

import (
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	gcpcloudstorage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestGcpSvcAccToken(t *testing.T) {

	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	gcpcloudstorage.GetSingleton().UpdateLoadK8sConfigMap("gcp-cloud-storage-config", "Dev", "")
	t.Run("TestGcpSvcAccToken", func(t *testing.T) {
		gact := GetSingleton()
		gact.UpdateLoadK8sConfigMap("svc-acc-config", "Dev", "")

	})

}
