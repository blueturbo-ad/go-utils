package k8s_informer

import (
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	logger "github.com/blueturbo-ad/go-utils/zap_loggerex"
	"github.com/stretchr/testify/assert"
)

func TestK8sInformer(t *testing.T) {
	t.Run("xxx", func(t *testing.T) {
		environment.Init()
		if err := GetInformerSingleton().SetUp(); err != nil {
			assert.NoError(t, err)
		}
		GetInformerSingleton().RegisterCacheInitFun("dsp-logger", logger.GetSingleton().UpdateLoadK8sConfigMap)
		GetInformerSingleton().Run()
	})
}
