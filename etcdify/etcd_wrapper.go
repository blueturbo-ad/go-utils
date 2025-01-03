package etcdify

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
etcd 订阅处理

	WatchKey 为订阅函数 通过传入的 key 进行订阅  callback为回调函数
*/

var (
	instance *EtcdWatcher
	once     sync.Once
)

type EtcdWatcher struct {
	Client  *clientv3.Client
	Timeout time.Duration
}

func GetSingleton() *EtcdWatcher {
	once.Do(func() {
		instance = new(EtcdWatcher)
	})
	return instance
}

func GetEtcder() *EtcdWatcher {
	return GetSingleton()
}

// 第一次调用接口  以后使用GetEtcder获取对象
func NewWatcher(confPath string, env string) error {
	var e = new(config_manage.EtcdifyConfig)
	err := e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}

	GetSingleton().Timeout = time.Duration(e.Timeout)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Hostname,
		DialTimeout: time.Duration(e.Timeout),
	})
	if err != nil {
		return fmt.Errorf("failed to create etcd client: %w", err)
	}

	GetSingleton().Client = cli
	return nil
}

func (w *EtcdWatcher) WatchKey(env string, ctx context.Context, key string, callback func(env string, eventType string, key string, value string)) {
	go func() {
		watchChan := w.Client.Watch(ctx, key)
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				callback(env, event.Type.String(), string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}()
}

func (w *EtcdWatcher) Close() error {
	return w.Client.Close()
}
