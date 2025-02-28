package k8s_informer

import (
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	logger "github.com/blueturbo-ad/go-utils/zap_loggerex"
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
		GetInformerSingleton().RegisterCacheInitFun("dsp-logger", logger.GetSingleton().UpdateLoadK8sConfigMap)
		GetInformerSingleton().Run()
	})
}
