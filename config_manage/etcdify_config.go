package config_manage

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

/* LoggerConfig 日志配置
 * Name 用于表示日志的名称
 */

type EtcdifyConfig struct {
	Hostname []string `yaml:"hostname"`
	Ip       string   `yaml:"ip"`
	Port     string   `yaml:"port"`
	Timeout  int64    `yaml:"timeout"`
}

func (l *EtcdifyConfig) LoadConfig(filePath string, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadFileConfig(filePath, env)
	if err != nil {
		return err
	}
	if info == nil {
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

func (l *EtcdifyConfig) LoadMemoryConfig(buf []byte, env string) error {
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
	err = yaml.Unmarshal(data, &l)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
