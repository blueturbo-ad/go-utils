package feishu

import (
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

type FeishuManage struct {
	Config [2]*config_manage.FeishuConfig
	index  int
}

var (
	instance *FeishuManage
	once     sync.Once
)

func GetInstance() *FeishuManage {
	once.Do(func() {
		instance = &FeishuManage{
			Config: [2]*config_manage.FeishuConfig{new(config_manage.FeishuConfig), new(config_manage.FeishuConfig)},
			index:  -1,
		}
	})
	return instance
}

// 函数用于内存更新etcd配置
func (l *FeishuManage) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "logger":
		var e = new(config_manage.FeishuConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		return l.UpdateLogger(e)
	default:
		return fmt.Errorf("unknown UpdateFromEtcd: %s", key)
	}
}

func (f *FeishuManage) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.FeishuConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}

	return f.UpdateLogger(e)
}

func (l *FeishuManage) UpdateLogger(config *config_manage.FeishuConfig) error {

	return nil
}
