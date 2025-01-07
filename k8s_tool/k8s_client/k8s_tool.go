package k8sclient

import (
	"os"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	instance    *K8sTool
	once        sync.Once
	EmptyString = ""
)

type K8sTool struct {
	client *kubernetes.Clientset
}

func GetSingleton() *K8sTool {
	once.Do(func() {
		instance = new(K8sTool)
	})
	return instance
}

func (k *K8sTool) SetUp() error {
	var config *rest.Config
	var err error
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		// 在集群内运行，使用 InClusterConfig
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	} else {
		// 在集群外运行，使用 kubeconfig 文件
		var kubeconfig string
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			kubeconfig = os.Getenv("KUBECONFIG")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k.client = client
	return nil
}

func (k *K8sTool) GetClient() *kubernetes.Clientset {
	return k.client
}
