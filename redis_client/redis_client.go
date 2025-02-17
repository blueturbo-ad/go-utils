package redisclient

import (
	"fmt"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
	redisconfigmanger "github.com/blueturbo-ad/go-utils/config_manage"

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
	return l.refreshRedisClient(e)
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
		if err := r.refreshRedisClient(e); err != nil {
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
	return f.refreshRedisClient(e)
}
func (r *RedisClientManager) refreshRedisClient(confs *redisconfigmanger.RedisConfigManager) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	newIndex := (r.index + 1) % 2
	r.ReadClient[newIndex] = make(map[string]*redis.ClusterClient)
	r.WriteClient[newIndex] = make(map[string]*redis.ClusterClient)
	for _, v := range *confs.Config {
		r.ReadClient[newIndex][v.Name] = r.BuildReadRedisClient(&v)
		r.WriteClient[newIndex][v.Name] = r.BuildWriteRedisClient(&v)
	}
	r.index = newIndex
	return nil
}

func (r *RedisClientManager) BuildWriteRedisClient(conf *redisconfigmanger.RedisConfig) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           []string{conf.WritePool.Nodes[0]},
		Password:        conf.WritePool.Password,
		ReadTimeout:     time.Duration(conf.WritePool.Timeout) * time.Millisecond,
		WriteTimeout:    time.Duration(conf.WritePool.Timeout) * time.Millisecond,
		PoolSize:        conf.WritePool.PoolSize,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 30 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.DB = conf.WritePool.Database
			opt.Password = conf.ReadPool.Password
			return redis.NewClient(opt)
		},
	})
}

func (r *RedisClientManager) BuildReadRedisClient(conf *redisconfigmanger.RedisConfig) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           []string{conf.ReadPool.Nodes[0]},
		Password:        conf.ReadPool.Password,
		ReadTimeout:     time.Duration(conf.ReadPool.Timeout) * time.Millisecond,
		WriteTimeout:    time.Duration(conf.ReadPool.Timeout) * time.Millisecond,
		PoolSize:        conf.ReadPool.PoolSize,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 30 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.DB = conf.ReadPool.Database
			opt.Password = conf.ReadPool.Password
			return redis.NewClient(opt)
		},
	})
}
