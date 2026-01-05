package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

/*
curused: Dev
Dev:

	version: 1
	postgresql_conf:
	  - name: mis_postgresql
	    write_pool:
	      database: blueturbo_business
	      pool_size: 20
	      timeout: 5000
	      port: 5432
	      nodes:
	        - 34.142.238.89
	      user: blueturbo
	      password: 7n<{qqu2pgh97ns~
	    read_pool:
	      write_pool:
	      database: blueturbo_business
	      pool_size: 20
	      timeout: 5000
	      port: 5432
	      nodes:
	        - 34.142.238.89
	      user: blueturbo_read
	      password: Cu5B3%l3M1pRXa0q

Pro:

	version: 1
	postgresql_conf:
	  - name: mis_postgresql
	    write_pool:
	      database: blueturbo_business
	      pool_size: 20
	      timeout: 5000
	      port: 5432
	      nodes:
	        - 34.142.238.89
	      user: blueturbo
	      password: 7n<{qqu2pgh97ns~
	    read_pool:
	      write_pool:
	      database: blueturbo_business
	      pool_size: 20
	      timeout: 5000
	      port: 5432
	      nodes:
	        - 34.142.238.89
	      user: blueturbo_read
	      password: Cu5B3%l3M1pRXa0q
*/
type PostgreSQLConfigManager struct {
	Config  *[]PostgreSQLConfig `yaml:"postgresql_conf"`
	Version string              `yaml:"version"`
}
type PostgreSQLConfig struct {
	Name      string `yaml:"name"`
	WritePool struct {
		Database string   `yaml:"database"`
		PoolSize int      `yaml:"pool_size"`
		Timeout  int      `yaml:"timeout"`
		Port     int      `yaml:"port"`
		Nodes    []string `yaml:"nodes"`
		User     string   `yaml:"user"`
		Password string   `yaml:"password"`
	} `yaml:"write_pool"`
	ReadPool struct {
		Database string   `yaml:"database"`
		PoolSize int      `yaml:"pool_size"`
		Timeout  int      `yaml:"timeout"`
		Port     int      `yaml:"port"`
		Nodes    []string `yaml:"nodes"`
		User     string   `yaml:"user"`
		Password string   `yaml:"password"`
	} `yaml:"read_pool"`
}

func (p *PostgreSQLConfigManager) LoadK8sConfigMap(configMapName, env string) error {
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
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
func (p *PostgreSQLConfigManager) LoadConfig(filePath string, env string) error {
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
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *PostgreSQLConfigManager) LoadMemoryConfig(buf []byte, env string) error {
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
