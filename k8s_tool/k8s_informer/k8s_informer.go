package k8s_informer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	loggerex "github.com/blueturbo-ad/go-utils/zap_loggerex"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	Informerinstance *Informer
	Informeronce     sync.Once
)

func GetInformerSingleton() *Informer {
	Informeronce.Do(func() {
		Informerinstance = new(Informer)
	})
	return Informerinstance
}

type Informer struct {
	cacheInitFuns map[string]func(configMapName, env string) error
	k8sClient     *kubernetes.Clientset
}

func (i *Informer) RegisterCacheInitFun(key string, fun func(configMapName, env string) error) {
	i.cacheInitFuns[key] = fun
}

func (i *Informer) SetUp() error {
	var err error
	i.cacheInitFuns = make(map[string]func(configMapName, env string) error)
	i.k8sClient = k8sclient.GetSingleton().GetClient()
	if err != nil {
		return err
	}
	return nil
}

func (i *Informer) Run() {
	// 创建 Informer 工厂
	namespace := environment.GetSingleton().GetNamespace()
	loggerex.GetSingleton().Info("system_logger", "namespace: %s", namespace)
	factory := informers.NewSharedInformerFactoryWithOptions(i.k8sClient, time.Minute*10, informers.WithNamespace(namespace))

	// 创建 ConfigMap Informer
	informer := factory.Core().V1().ConfigMaps().Informer()
	// 添加事件处理程序
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			env := environment.GetSingleton().GetEnv()
			loggerex.GetSingleton().Info("system_logger", "add config map: %s", configMap.Name)
			if initFunc, exists := i.cacheInitFuns[configMap.Name]; exists {
				if err := initFunc(configMap.Name, env); err != nil {
					loggerex.GetSingleton().Error("system_logger", "add config map error : %s", err.Error())
				}
			} else {
				loggerex.GetSingleton().Error("system_logger", "No init function found for config map: %s", configMap.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newConfigMap := newObj.(*corev1.ConfigMap)
			env := environment.GetSingleton().GetEnv()
			loggerex.GetSingleton().Info("system_logger", "add config map: %s \n", newConfigMap.Name)
			if initFunc, exists := i.cacheInitFuns[newConfigMap.Name]; exists {
				if err := initFunc(newConfigMap.Name, env); err != nil {
					loggerex.GetSingleton().Error("system_logger", "update config map error : %s \n", err.Error())
				}
			} else {
				loggerex.GetSingleton().Error("system_logger", "No init function found for config map: %s", newConfigMap.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			//TODO
			fmt.Println("delete config map: ", configMap.Name)
		},
	})
	// 启动 Informer
	stopCh := make(chan struct{})
	defer close(stopCh)

	go informer.Run(stopCh)
	// 等待缓存同步
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		loggerex.GetSingleton().Error("system_logger", "Error waiting for cache to sync")
	}
	// 等待信号以退出程序
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
