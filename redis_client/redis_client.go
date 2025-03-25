package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
	redisconfigmanger "github.com/blueturbo-ad/go-utils/config_manage"
	gcpcloudstorage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"

	"github.com/redis/go-redis/v9"
)

var (
	instance    *RedisClientManager
	once        sync.Once
	EmptyString = ""
)

// 双缓存管理器
type RedisClientManager struct {
	ReadClient  [2]map[string]*redis.ClusterClient
	WriteClient [2]map[string]*redis.ClusterClient
	index       int
	rwMutex     sync.RWMutex
}

func GetSingleton() *RedisClientManager {
	once.Do(func() {
		instance = &RedisClientManager{
			index: -1,
		}

	})
	return instance
}

func (r *RedisClientManager) GetReadClient(name string) *redis.ClusterClient {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()
	if r.ReadClient[r.index][name] != nil {
		return r.ReadClient[r.index][name]
	}
	return nil
}

func (r *RedisClientManager) GetWriteClient(name string) *redis.ClusterClient {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()
	if r.WriteClient[r.index][name] != nil {
		return r.WriteClient[r.index][name]
	}
	return nil
}

func (l *RedisClientManager) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.RedisConfigManager)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return l.refreshRedisClient(e, env)
}

// 函数用于内存更新etcd配置
func (r *RedisClientManager) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "redis-conf":
		var e = new(config_manage.RedisConfigManager)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := r.refreshRedisClient(e, env); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (f *RedisClientManager) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.RedisConfigManager)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return f.refreshRedisClient(e, env)
}
func (r *RedisClientManager) refreshRedisClient(confs *redisconfigmanger.RedisConfigManager, env string) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	newIndex := (r.index + 1) % 2
	r.ReadClient[newIndex] = make(map[string]*redis.ClusterClient)
	r.WriteClient[newIndex] = make(map[string]*redis.ClusterClient)
	for _, v := range *confs.Config {
		r.ReadClient[newIndex][v.Name] = r.BuildReadRedisClient(&v, env)
		r.WriteClient[newIndex][v.Name] = r.BuildWriteRedisClient(&v, env)
	}
	r.index = newIndex
	return nil
}

func (r *RedisClientManager) BuildWriteRedisClient(conf *redisconfigmanger.RedisConfig, env string) *redis.ClusterClient {
	// username, password, _ := r.retrieveTokenFunc(conf)
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           []string{conf.WritePool.Nodes[0]},
		ReadTimeout:     time.Duration(conf.WritePool.Timeout) * time.Millisecond,
		WriteTimeout:    time.Duration(conf.WritePool.Timeout) * time.Millisecond,
		PoolSize:        conf.WritePool.PoolSize,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 30 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
}

func (r *RedisClientManager) BuildReadRedisClient(conf *redisconfigmanger.RedisConfig, env string) *redis.ClusterClient {
	// username, password, _ := r.retrieveTokenFunc(conf)
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           []string{conf.ReadPool.Nodes[0]},
		ReadTimeout:     time.Duration(conf.ReadPool.Timeout) * time.Millisecond,
		WriteTimeout:    time.Duration(conf.ReadPool.Timeout) * time.Millisecond,
		PoolSize:        conf.ReadPool.PoolSize,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 30 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
}

func (r *RedisClientManager) retrieveTokenFunc(conf *redisconfigmanger.RedisConfig, env string) (string, string, error) {
	if env == "Dev" {
		return "", "", nil
	}
	ctx := context.Background()
	client := gcpcloudstorage.GetSingleton().GetClient("dsp_bucket")
	if client == nil {
		return EmptyString, EmptyString, fmt.Errorf("cloud storage client is nil")
	}
	wc, err := client.Object("account_token/access_token.json").NewReader(ctx)
	if err != nil {
		return EmptyString, EmptyString, err
	}
	defer wc.Close()
	// 读取对象内容
	data, err := io.ReadAll(wc)
	if err != nil {
		return EmptyString, EmptyString, err
	}
	fmt.Println(string(data))
	// 反序列化为 map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return EmptyString, EmptyString, err
	}
	var token string
	if result[conf.Name] != nil {
		if val, ok := result[conf.Name].(string); ok {
			token = val
		}
	}
	username := "default"
	password := token
	fmt.Println("username:", username, "password:", password)
	return username, password, nil
}
