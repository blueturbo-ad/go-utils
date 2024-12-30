package etcdify

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
etcd 订阅处理

	WatchKey 为订阅函数 通过传入的 key 进行订阅  callback为回调函数
*/
type EtcdWatcher struct {
	client  *clientv3.Client
	timeout time.Duration
}

func NewWatcher(endpoints []string, timeout time.Duration) (*EtcdWatcher, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &EtcdWatcher{
		client:  cli,
		timeout: timeout,
	}, nil
}

func (w *EtcdWatcher) WatchKey(env string, ctx context.Context, key string, callback func(env string, eventType string, key string, value string)) {
	watchChan := w.client.Watch(ctx, key)
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				callback(env, event.Type.String(), string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}()
}

func (w *EtcdWatcher) Close() error {
	return w.client.Close()
}
