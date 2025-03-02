package k8s_informer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sync"
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
	Informer      *cache.SharedIndexInformer
	StartErrChan  chan error
	cacheFunc     map[string]bool // 这里只记录add 的状态 add 只在首次启动的时候会加载
	Ssucchan      chan bool       // 首次启动等待时间是30s 如果我们提前启动完成就给这个信号来通知
}

func (i *Informer) RegisterCacheInitFun(key string, fun func(configMapName, env string) error) {
	i.cacheInitFuns[key] = fun
	i.cacheFunc[key] = false
}

func (i *Informer) SetUp() error {
	var err error
	i.cacheInitFuns = make(map[string]func(configMapName, env string) error)
	i.k8sClient = k8sclient.GetSingleton().GetClient()
	namespace := environment.GetSingleton().GetNamespace()
	factory := informers.NewSharedInformerFactoryWithOptions(i.k8sClient, 60*time.Minute, informers.WithNamespace(namespace))
	// 创建 ConfigMap Informer
	informer := factory.Core().V1().ConfigMaps().Informer()
	i.Informer = &informer
	i.StartErrChan = make(chan error, 100)
	i.cacheFunc = make(map[string]bool)
	i.Ssucchan = make(chan bool, 10)

	if err != nil {
		return err
	}
	return nil
}

// 这里关闭的自动同步是在同步日志配置的时候 由于buff的切换丢失了上一次的file 对象
func (i *Informer) Run() {
	// 创建 Informer 工厂
	informer := (*i.Informer)
	// 添加事件处理程序
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap := obj.(*corev1.ConfigMap)
			env := environment.GetSingleton().GetEnv()
			loggerex.GetSingleton().Info("system_logger", "add config map: %s", configMap.Name)
			if initFunc, exists := i.cacheInitFuns[configMap.Name]; exists {
				if err := initFunc(configMap.Name, env); err != nil {
					msg := fmt.Sprintf("add config map error: %s", err.Error())
					loggerex.GetSingleton().Error("system_logger", "%s", msg)
					i.StartErrChan <- fmt.Errorf("%s", msg)
				} else {
					i.cacheFunc[configMap.Name] = true
				}
			} else {
				msg := fmt.Sprintf("No init function found for config map: %s", configMap.Name)
				loggerex.GetSingleton().Warn("system_logger", "%s", msg)
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
					msg := fmt.Sprintf("update config map error: %s", err.Error())
					loggerex.GetSingleton().Error("system_logger", "%s", msg)
				}
			} else {
				msg := fmt.Sprintf("No init function found for config map: %s", newConfigMap.Name)
				loggerex.GetSingleton().Warn("system_logger", "%s", msg)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 等待缓存同步
	go func() {
		for {
			if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
				i.StartErrChan <- errors.New("缓存同步失败")
			}
			select {
			case <-ctx.Done():
				fmt.Println("cache sync done")
				if err := i.CheckIsRun(); err != nil {
					i.StartErrChan <- err
					return
				}
				return
			default:
				if err := i.CheckIsRun(); err == nil {
					i.Ssucchan <- true
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	// 等待信号以退出程序
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func (i *Informer) CheckIsRun() error {
	for key, value := range i.cacheFunc {
		if !value {
			return fmt.Errorf("cache %s is not run", key)
		}
	}
	return nil
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
