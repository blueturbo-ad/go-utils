package rtaattr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	dsp_base_config "github.com/blueturbo-ad/go-utils/dsp_base_config"
	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestRtaAttrConf(t *testing.T) {

	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Dir(dir)
	p := dir + "/config/rta_attr_conf.yaml"
	t.Run("RtaAttrConf", func(t *testing.T) {
		// 读取yaml 文件里面的内容
		openFile, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		defer openFile.Close()
		fileInfo, err := openFile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSize := fileInfo.Size()
		data := make([]byte, fileSize)
		_, err = openFile.Read(data)
		if err != nil {
			log.Fatal(err)
		}
		dsp_base_config.GetSingleton().RegistHookFunc(GetSingleton().Reload, "rta_attr_conf")
		if err := dsp_base_config.GetSingleton().LoadFileConfig(p, environment.GetEnv(), "rta_attr_conf"); err != nil {
			fmt.Printf("Failed to load config file: %v\n", err)
		}
		t.Logf("RtaAttrConf: %+v", GetSingleton().GetRtaAttrConf("tiktok"))

	})

}
