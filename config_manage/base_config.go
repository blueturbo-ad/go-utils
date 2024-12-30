package config_manage

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

/* Config 用于表示整个配置文件
 * CurUsed 用于表示当前使用的环境
 * Dev 用于表示开发环境的配置
 */

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

func (c *ManagerConfig) LoadConfig(filePath string, env string) (*any, error) {
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
	//fmt.Println(fmt.Sprintf("env = %s CurUsed = %s", env, c.CurUsed))
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
