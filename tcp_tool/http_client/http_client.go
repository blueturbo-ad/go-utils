package httpclient

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

var (
	instance    *HttpClientManager
	once        sync.Once
	EmptyString = ""
)

type HttpClientManager struct {
	HttpClient [2]map[string]*http.Client

	index   int
	rwMutex sync.RWMutex
}

func GetSingleton() *HttpClientManager {
	once.Do(func() {
		instance = &HttpClientManager{
			index: -1,
		}

	})
	return instance
}

func (h *HttpClientManager) GetClient(name string) *http.Client {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()
	if h.HttpClient[h.index][name] != nil {
		return h.HttpClient[h.index][name]
	}
	return nil
}

func (h *HttpClientManager) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.HttpClientConfig)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return h.BuildClient(e)
}

// 函数用于内存更新etcd配置
func (h *HttpClientManager) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "redis-conf":
		var e = new(config_manage.HttpClientConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := h.BuildClient(e); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (h *HttpClientManager) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.HttpClientConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return h.BuildClient(e)
}

func (h *HttpClientManager) BuildClient(e *config_manage.HttpClientConfig) error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	h.index = (h.index + 1) % 2
	h.HttpClient[h.index] = make(map[string]*http.Client)

	for _, v := range *e.Config {
		client := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(v.TimeOut) * time.Millisecond,
					KeepAlive: time.Duration(v.KeepALive) * time.Millisecond,
				}).DialContext,
				MaxIdleConns:        v.MaxIdleConns,
				MaxIdleConnsPerHost: v.MaxIdleConnsPerHost,
				IdleConnTimeout:     time.Duration(v.IdleConnTimeout) * time.Millisecond,
			},
		}
		h.HttpClient[h.index][v.Name] = &client
	}
	return nil
}
