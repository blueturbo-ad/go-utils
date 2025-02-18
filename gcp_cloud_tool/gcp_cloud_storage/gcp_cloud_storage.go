package gcpcloudstorage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/blueturbo-ad/go-utils/config_manage"
)

var (
	instance    *GcpCloudStorage
	once        sync.Once
	EmptyString = ""
)

type GcpCloudStorage struct {
	Clients map[string]*storage.BucketHandle
}

func GetSingleton() *GcpCloudStorage {
	once.Do(func() {
		instance = &GcpCloudStorage{}
	})
	return instance
}

func (g *GcpCloudStorage) GetClient(bucketName string) *storage.BucketHandle {
	return g.Clients[bucketName]
}

func (g *GcpCloudStorage) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.GcpCloudStorageTokenConfig)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return g.BucketClient(e.GcpAttr)
}

// 函数用于内存更新etcd配置
func (r *GcpCloudStorage) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "redis-conf":
		var e = new(config_manage.GcpCloudStorageTokenConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := r.BucketClient(e.GcpAttr); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (f *GcpCloudStorage) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.GcpCloudStorageTokenConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return f.BucketClient(e.GcpAttr)
}

func (g *GcpCloudStorage) BucketClient(confs []config_manage.CloudStorage) error {
	for _, conf := range confs {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()
		// Creates a Bucket instance.
		bucket := client.Bucket(conf.Bucket)
		if g.Clients == nil {
			g.Clients = make(map[string]*storage.BucketHandle)
		}
		g.Clients[conf.Bucket] = bucket
	}

	return nil
}
