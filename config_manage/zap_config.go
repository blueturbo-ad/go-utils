package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

/* LoggerConfig 日志配置
 * Name 用于表示日志的名称
 */

type LoggerConfig struct {
	Name       string `yaml:"name"`
	Info       string `yaml:"info"`
	Error      string `yaml:"error"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	Level      int    `yaml:"level"`
	Async      bool   `yaml:"async"`
}

// Config 用于表示整个配置文件
type ZapLoggerConfig struct {
	Loggers []LoggerConfig `yaml:"zap_loggers"`
	Version string         `yaml:"version"`
}

func (l *ZapLoggerConfig) LoadK8sConfigMap(configMapName, env string) error {
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

	err = yaml.Unmarshal(data, &l)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *ZapLoggerConfig) LoadZapConfig(filePath string, env string) error {
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
	err = yaml.Unmarshal(data, &l)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *ZapLoggerConfig) LoadMemoryZapConfig(buf []byte, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadMemoryConfig(buf, env)
	if err != nil {
		return err
	}
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &l)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
