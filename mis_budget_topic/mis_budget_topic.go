package misbudgettopic

import (
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

var (
	instance    *MisBudgetTopicManager
	once        sync.Once
	EmptyString = ""
)

// 双缓存管理器
type MisBudgetTopicManager struct {
	Config map[string]*config_manage.MisBudgetTopicConfig
}

func GetSingleton() *MisBudgetTopicManager {
	once.Do(func() {
		instance = &MisBudgetTopicManager{}

	})
	return instance
}

func (m *MisBudgetTopicManager) GetConfig(name string) *config_manage.MisBudgetTopicConfig {
	m.Config[name].Name = name
	return m.Config[name]
}

func (m *MisBudgetTopicManager) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.MisBudgetTopicConfigManager)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return m.refreshConfig(e, env)
}

func (m *MisBudgetTopicManager) UpdateLoadFileConfig(filePath, env string) error {
	var e = new(config_manage.MisBudgetTopicConfigManager)
	err := e.LoadConfig(filePath, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return m.refreshConfig(e, env)
}

func (m *MisBudgetTopicManager) refreshConfig(e *config_manage.MisBudgetTopicConfigManager, env string) error {
	m.Config = make(map[string]*config_manage.MisBudgetTopicConfig)
	for _, v := range *e.Config {
		m.Config[v.Name] = &v
	}
	return nil
}
