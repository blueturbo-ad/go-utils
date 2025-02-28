package k8s_informer

import (
	"context"
	"errors"
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
	Informer      cache.SharedIndexInformer
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

// 这里关闭的自动同步是在同步日志配置的时候 由于buff的切换丢失了上一次的file 对象
func (i *Informer) Run() {
	// 创建 Informer 工厂
	namespace := environment.GetSingleton().GetNamespace()
	loggerex.GetSingleton().Info("system_logger", "namespace: %s", namespace)
	factory := informers.NewSharedInformerFactoryWithOptions(i.k8sClient, 60*time.Minute, informers.WithNamespace(namespace))

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
			oldConfigMap := oldObj.(*corev1.ConfigMap)
			env := environment.GetSingleton().GetEnv()
			loggerex.GetSingleton().Info("system_logger", "update config map: %s \n", newConfigMap.Name)
			if initFunc, exists := i.cacheInitFuns[newConfigMap.Name]; exists {
				if IsConfigMapEqual(oldConfigMap, newConfigMap) {
					return
				}
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

	i.Informer = informer

	go informer.Run(stopCh)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 等待缓存同步
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		select {
		case <-ctx.Done():
			// 检查是否因超时失败
			if ctx.Err() == context.DeadlineExceeded {
				panic(errors.New("同步超时：无法在30秒内完成缓存同步"))
			}
			panic(fmt.Errorf("同步被取消: %v", ctx.Err()))
		default:
			panic(errors.New("缓存同步失败"))
		}
	}
	// 等待信号以退出程序
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

// IsConfigMapEqual 判断两个 ConfigMap 是否相同
func IsConfigMapEqual(oldConfigMap, newConfigMap *corev1.ConfigMap) bool {
	// 比较元数据
	if oldConfigMap.Name != newConfigMap.Name ||
		oldConfigMap.Namespace != newConfigMap.Namespace ||
		oldConfigMap.ResourceVersion != newConfigMap.ResourceVersion {
		return false
	}

	// 比较数据字段
	if len(oldConfigMap.Data) != len(newConfigMap.Data) {
		return false
	}
	for key, oldValue := range oldConfigMap.Data {
		if newValue, ok := newConfigMap.Data[key]; !ok || newValue != oldValue {
			return false
		}
	}

	// 比较二进制数据字段
	if len(oldConfigMap.BinaryData) != len(newConfigMap.BinaryData) {
		return false
	}
	for key, oldValue := range oldConfigMap.BinaryData {
		if newValue, ok := newConfigMap.BinaryData[key]; !ok || !EqualByteSlices(newValue, oldValue) {
			return false
		}
	}

	return true
}

func EqualByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
