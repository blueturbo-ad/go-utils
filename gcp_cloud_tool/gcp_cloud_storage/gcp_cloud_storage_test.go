package gcpcloudstorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestReadGCPCloudStorageFile(t *testing.T) {
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	GetSingleton().UpdateLoadK8sConfigMap("gcp-cloud-storage-config", environment.GetEnv(), "")

	client := GetSingleton().GetClient("dsp_bucket")

	if client == nil {
		t.Errorf("GetClient error")
		return
	}

	reader, err := client.Object("ad_bid_server_data/country_label_id.json.stat").NewReader(context.Background())
	if err != nil {
		t.Errorf("ReadGCPCloudStorageFile error: %v", err)
		return
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("ReadGCPCloudStorageFile error: %v", err)
		return
	}

	fmt.Println(string(data))
}
