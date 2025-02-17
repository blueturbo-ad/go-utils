package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Name      string `yaml:"name"`
	WritePool struct {
		Database int      `yaml:"database"`
		PoolSize int      `yaml:"pool_size"`
		Password string   `yaml:"password"`
		Timeout  int      `yaml:"timeout"`
		Nodes    []string `yaml:"nodes"`
	} `yaml:"write_pool"`
	ReadPool struct {
		Database int      `yaml:"database"`
		PoolSize int      `yaml:"pool_size"`
		Password string   `yaml:"password"`
		Timeout  int      `yaml:"timeout"`
		Nodes    []string `yaml:"nodes"`
	} `yaml:"read_pool"`
}

type RedisConfigManager struct {
	Config  *[]RedisConfig `yaml:"redis_conf"`
	Version string         `yaml:"version"`
}

func (r *RedisConfigManager) LoadK8sConfigMap(configMapName, env string) error {
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
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (r *RedisConfigManager) LoadConfig(filePath string, env string) error {
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
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *RedisConfigManager) LoadMemoryConfig(buf []byte, env string) error {
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
