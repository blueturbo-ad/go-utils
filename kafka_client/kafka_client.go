package kafkaclient

import (
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	instance    *KafkaClientManager
	once        sync.Once
	EmptyString = ""
)

type KafkaClientManager struct {
	ProducerClient [2]map[string]*kafka.Producer
	ConsumerClient [2]map[string]*kafka.Consumer
	index          int
	rwMutex        sync.RWMutex
}

func GetSingleton() *KafkaClientManager {
	once.Do(func() {
		instance = NewKafkaClientManager()

	})
	return instance
}

func NewKafkaClientManager() *KafkaClientManager {
	return &KafkaClientManager{
		ProducerClient: [2]map[string]*kafka.Producer{},
		ConsumerClient: [2]map[string]*kafka.Consumer{},
		index:          -1,
	}
}

func (k *KafkaClientManager) GetProducerClient(name string) *kafka.Producer {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	if k.ProducerClient[k.index][name] != nil {
		return k.ProducerClient[k.index][name]
	}
	return nil
}

func (k *KafkaClientManager) GetConsumerClient(name string) *kafka.Consumer {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	if k.ConsumerClient[k.index][name] != nil {
		return k.ConsumerClient[k.index][name]
	}
	return nil
}

func (k *KafkaClientManager) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.KafkaConfigManage)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("kafka client  LoadK8sConfigMap is error %s", err.Error())
	}
	return k.buildKafkaClient(e)
}

// 函数用于内存更新etcd配置
func (k *KafkaClientManager) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "kafka-config":
		var e = new(config_manage.KafkaConfigManage)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := k.buildKafkaClient(e); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (k *KafkaClientManager) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.KafkaConfigManage)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return k.buildKafkaClient(e)
}

func (k *KafkaClientManager) buildKafkaClient(e *config_manage.KafkaConfigManage) error {
	for _, v := range *e.Config {
		p, err := k.buildProducer(&v)
		if err != nil {
			return err
		}
		c, err := k.buildConsumer(&v)
		if err != nil {
			return err
		}
		k.rwMutex.Lock()
		k.index = (k.index + 1) % 2
		if k.ProducerClient[k.index] == nil {
			k.ProducerClient[k.index] = make(map[string]*kafka.Producer)
		}
		if k.ConsumerClient[k.index] == nil {
			k.ConsumerClient[k.index] = make(map[string]*kafka.Consumer)
		}
		k.ProducerClient[k.index][v.Name] = p
		k.ConsumerClient[k.index][v.Name] = c

	}
	return nil
}

func (k *KafkaClientManager) buildProducer(conf *config_manage.KafkaConfig) (*kafka.Producer, error) {
	// 创建生产者配置
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Producer.Broker,
		"client.id":         conf.Producer.Producer,
		"acks":              "all",
		"security.protocol": conf.Producer.Protocol,
		"sasl.mechanism":    conf.Producer.Mechanism,
		"sasl.username":     conf.Producer.Username,
		"sasl.password":     conf.Producer.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %s", err)
	}

	return p, nil

}

func (k *KafkaClientManager) buildConsumer(conf *config_manage.KafkaConfig) (*kafka.Consumer, error) {
	// 创建消费者配置
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Customer.Broker,
		"group.id":          conf.Customer.Group,
		"auto.offset.reset": conf.Customer.Reset,
		"security.protocol": conf.Customer.Protocol,
		"sasl.mechanism":    conf.Customer.Mechanism,
		"sasl.username":     conf.Customer.Username,
		"sasl.password":     conf.Customer.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}
	return c, nil
}
