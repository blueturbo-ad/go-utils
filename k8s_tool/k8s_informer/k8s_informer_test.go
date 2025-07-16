package k8s_informer

import (
	"os"
	"testing"

	dsp_base_config "github.com/blueturbo-ad/go-utils/dsp_base_config"
	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	"github.com/stretchr/testify/assert"
)

func TestK8sInformer(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("xxx", func(t *testing.T) {
		environment.Init()
		if err := GetInformerSingleton().SetUp(); err != nil {
			assert.NoError(t, err)
		}
		Informer := GetInformerSingleton().Informer
		assert.NotNil(t, Informer)
		GetInformerSingleton().RegisterCacheInitFun("bid-server-iconf-config", func(configMapName string, env string, hookName string) error {
			// Provide a default value for hookName or adjust as needed
			return dsp_base_config.GetSingleton().LoadK8sConfigMap(configMapName, env, hookName)
		})
		// GetInformerSingleton().RegisterCacheInitFun("feishu", logger.GetSingleton().UpdateLoadK8sConfigMap)

		GetInformerSingleton().Run()
		for {
			select {
			case err := <-GetInformerSingleton().StartErrChan:
				assert.NoError(t, err)
			case <-GetInformerSingleton().Ssucchan:
				break
			}
		}
	})
}
