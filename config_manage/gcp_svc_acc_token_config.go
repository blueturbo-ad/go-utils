package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

type CloudAcc struct {
	Name         string `yaml:"name"`
	AccountEmail string `yaml:"account_email"`
}

type GcpSvcAccountTokenConfig struct {
	CloudAcc []CloudAcc `yaml:"cloud_acc"`
	Version  string     `yaml:"version"`
}

func (g *GcpSvcAccountTokenConfig) LoadK8sConfigMap(configMapName, env string) error {
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
	err = yaml.Unmarshal(data, &g)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (r *GcpSvcAccountTokenConfig) LoadConfig(filePath string, env string) error {
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

func (g *GcpSvcAccountTokenConfig) LoadMemoryConfig(buf []byte, env string) error {
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
	err = yaml.Unmarshal(data, &g)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
