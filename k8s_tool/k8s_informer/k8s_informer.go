package k8s_informer

import (
	"fmt"
	"sync"
	"time"

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
	CacheInitFuns map[string]func() error
	k8sClient     *kubernetes.Clientset
}

func (i *Informer) RegisterCacheInitFun(key string, fun func() error) {
	i.CacheInitFuns[key] = fun
}

func (i *Informer) SetUp() error {
	var err error
	i.k8sClient = k8sclient.GetSingleton().GetClient()
	if err != nil {
		return err
	}
	return nil
}

func (i *Informer) Run() {
	// 创建 Informer 工厂
	factory := informers.NewSharedInformerFactoryWithOptions(i.k8sClient, time.Minute*10, informers.WithNamespace("default"))

	// 创建 ConfigMap Informer
	informer := factory.Core().V1().ConfigMaps().Informer()
	// 添加事件处理程序
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			if err := i.CacheInitFuns[configMap.Name](); err != nil {
				loggerex.GetSingleton().Error("framework_logger", "add config map error : %s", err.Error())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newConfigMap := newObj.(*corev1.ConfigMap)
			if err := i.CacheInitFuns[newConfigMap.Name](); err != nil {
				loggerex.GetSingleton().Error("framework_logger", "add config map error : %s", err.Error())
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
		loggerex.GetSingleton().Error("framework_logger", "Error waiting for cache to sync")
	}
}