package bigtableclient

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/bigtable"
	"github.com/blueturbo-ad/go-utils/config_manage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	instance    *BigTableClient
	once        sync.Once
	EmptyString = ""
)

func GetSingleton() *BigTableClient {
	once.Do(func() {
		instance = NewBigTableClient()
	})
	return instance
}

type BigTableClient struct {
	Clients [2]*bigtable.Client
	index   int
	lock    sync.Mutex
}

func NewBigTableClient() *BigTableClient {
	return &BigTableClient{
		Clients: [2]*bigtable.Client{nil, nil},
	}
}

func (b *BigTableClient) GetClient() *bigtable.Client {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.Clients[b.index]
}

func (b *BigTableClient) UpdateLoadK8sConfigMap(configMapName, env string, hookName string) error {
	var e = new(config_manage.BigTableConfig)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return err
	}
	return b.CreateClient(e)
}

// 函数用于内存更新etcd配置
func (b *BigTableClient) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "logger":
		var e = new(config_manage.BigTableConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := b.CreateClient(e); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (b *BigTableClient) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.BigTableConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return b.CreateClient(e)
}

func (b *BigTableClient) CreateClient(e *config_manage.BigTableConfig) error {
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return err
	}

	client, err := bigtable.NewClient(ctx, e.Config.ProjectId, e.Config.InstanceId, option.WithCredentials(creds))
	if err != nil {
		return err
	}
	if err := b.UpdateLogger(client); err != nil {
		return err
	}
	return nil
}

func (b *BigTableClient) UpdateLogger(client *bigtable.Client) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.index = (b.index + 1) % 2
	b.Clients[b.index] = client
	return nil
}
