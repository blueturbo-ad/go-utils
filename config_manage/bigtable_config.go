package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

type TableConfig struct {
	ProjectId  string `yaml:"project_id"`
	InstanceId string `yaml:"instance_id"`
}

// Config 用于表示整个配置文件
type BigTableConfig struct {
	Config  *TableConfig `yaml:"zap_loggers"`
	Version string       `yaml:"version"`
}

func (l *BigTableConfig) LoadK8sConfigMap(configMapName, env string) error {
	var c = new(ManagerConfig)
	namespace := environment.GetSingleton().GetNamespace()
	info, err := c.LoadK8sConfigMap(namespace, configMapName, env)
	if err != nil {
		return err
	}
	if (*info) == nil {
		return fmt.Errorf("info is nil，")
	}
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &l.Config)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *BigTableConfig) LoadConfig(filePath string, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadFileConfig(filePath, env)
	if err != nil {
		return err
	}
	//fmt.Println(info)
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &l.Config)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *BigTableConfig) LoadMemoryConfig(buf []byte, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadMemoryConfig(buf, env)
	if err != nil {
		return err
	}
	// 解析 YAML 数据
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &l.Config)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
