package config_manage

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Url string `yaml:"url"`
}

// Config 用于表示整个配置文件
type FeishuConfig struct {
	Config  *Config `yaml:"zap_loggers"`
	Version string  `yaml:"version"`
}

func (l *FeishuConfig) LoadConfig(filePath string, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadConfig(filePath, env)
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

func (l *FeishuConfig) LoadMemoryConfig(buf []byte, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadMemoryConfig(buf, env)
	if err != nil {
		return err
	}
	fmt.Println(info)
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