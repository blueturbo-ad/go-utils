package config_manage

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type WritePoolConfig struct {
	DataBase int      `yaml:"database"`
	TimeOut  int      `yaml:"timeout"`
	PoolSize int      `yaml:"poolsize"`
	Nodes    []string `yaml:"nodes"`
}

type ReadPoolConfig struct {
	DataBase int      `yaml:"database"`
	TimeOut  int      `yaml:"timeout"`
	PoolSize int      `yaml:"poolsize"`
	Nodes    []string `yaml:"nodes"`
}

type RedisConfig struct {
	Name      string          `yaml:"name"`
	WritePool WritePoolConfig `yaml:"writepool"`
	ReadPool  ReadPoolConfig  `yaml:"readpool"`
}

type RedisConfigManager struct {
	Config *RedisConfig `yaml:"redis"`
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
