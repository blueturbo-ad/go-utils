package config_manage

import (
	"context"
	"fmt"
	"os"

	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/* Config 用于表示整个配置文件
 * CurUsed 用于表示当前使用的环境
 * Dev 用于表示开发环境的配置
 */
type ManagerConfigInterface interface {
	LoadFileConfig(filePath string, env string) (*any, error)
	LoadMemoryConfig(buf []byte, env string) (*any, error)
	LoadK8sConfigMap(env string) (*any, error)
}

type ManagerConfig struct {
	CurUsed string      ` yaml:"curused"`
	Dev     interface{} `yaml:"Dev"`
	Pro     interface{} `yaml:"Pro"`
	Pre     interface{} `yaml:"Pre"`
	Test    interface{} `yaml:"Test"`
}

const (
	ErrorEnvNotfound  = "failed to read file: %w"
	ErroryamlNotfound = "failed to parse YAML: %w"
	KeyFieldTag       = "yaml"
)

func (c *ManagerConfig) LoadK8sConfigMap(namespace, configMapName, env string) (*any, error) {
	// 读取 YAML 文件
	k8s_client := k8sclient.GetSingleton().GetClient()
	if k8s_client == nil {
		return nil, fmt.Errorf("k8s client is nil")
	}
	configMap, err := k8s_client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %v", err)
	}
	// val := configMap.Data
	// conf, ok := configMap.Data[configMapName]
	// if !ok {
	// 	return nil, fmt.Errorf("configmap %s not found", configMapName)
	// }
	var conf string
	for _, v := range configMap.Data {
		conf = v
		break
	}
	if conf == "" {
		return nil, fmt.Errorf("configmap %s not found", configMapName)
	}

	fmt.Printf("base data: %s", conf)
	err = yaml.Unmarshal([]byte(conf), &c)
	if err != nil {
		return nil, fmt.Errorf(ErroryamlNotfound, err)
	}
	return c.GetEnvironmentConfig(env)
}

func (c *ManagerConfig) LoadFileConfig(filePath string, env string) (*any, error) {
	// 读取 YAML 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(ErrorEnvNotfound, err)
	}

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf(ErroryamlNotfound, err)
	}

	return c.GetEnvironmentConfig(env)
}

func (c *ManagerConfig) LoadMemoryConfig(buf []byte, env string) (*any, error) {
	// 解析 YAML 数据
	err := yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf(ErroryamlNotfound, err)
	}

	return c.GetEnvironmentConfig(env)
}

func (c *ManagerConfig) GetEnvironmentConfig(env string) (*any, error) {
	if env == "" {
		env = c.CurUsed
	}
	switch env {
	case "Dev":
		return &c.Dev, nil
	case "Pro":
		return &c.Pro, nil
	case "Pre":
		return &c.Pro, nil
	case "Test":
		return &c.Pro, nil
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
}
