package kafkaclient

import (
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	instance    *KafkaClientManager
	once        sync.Once
	EmptyString = ""
)

type KafkaClientManager struct {
	Config         [2]map[string]*config_manage.KafkaConfig
	ConsumerClient [2]map[string]map[string]*kafka.Consumer
	ProducerClient [2]map[string]*kafka.Producer // 添加 Producer 缓存
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
		Config:         [2]map[string]*config_manage.KafkaConfig{},
		ConsumerClient: [2]map[string]map[string]*kafka.Consumer{},
		ProducerClient: [2]map[string]*kafka.Producer{},
		index:          -1,
	}
}

func (k *KafkaClientManager) GetProducerClient(name string) (*kafka.Producer, error) {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	// 检查是否已有缓存的 Producer
	if k.index >= 0 && k.ProducerClient[k.index] != nil {
		if producer, exists := k.ProducerClient[k.index][name]; exists && producer != nil {
			// 安全的访问
			return producer, nil
		}
	}
	// 如果没有缓存，创建新的 Producer 并缓存
	if k.Config[k.index][name] != nil {
		config := k.Config[k.index][name]
		producer, err := k.buildProducer(config)
		if err != nil {
			return nil, err
		}

		// 缓存 Producer
		if k.ProducerClient[k.index] == nil {
			k.ProducerClient[k.index] = make(map[string]*kafka.Producer)
		}
		k.ProducerClient[k.index][name] = producer

		return producer, nil
	}
	return nil, fmt.Errorf("kafka client GetProducerClient is error")
}

func (k *KafkaClientManager) GetConsumerClient(name string, group string) (*kafka.Consumer, error) {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	if consumer, exists := k.ConsumerClient[k.index][name]; exists && consumer != nil {
		if c, ok := consumer[group]; ok {
			return c, nil
		}
	}
	// 如果没有缓存，创建新的 Producer 并缓存
	if k.Config[k.index][name] != nil {
		config := k.Config[k.index][name]
		consumer, err := k.buildConsumer(config, group)
		if err != nil {
			return nil, err
		}

		// 缓存 Producer
		// 正确初始化 ConsumerClient
		if k.ConsumerClient[k.index] == nil {
			k.ConsumerClient[k.index] = make(map[string]map[string]*kafka.Consumer)
		}
		if k.ConsumerClient[k.index][name] == nil {
			k.ConsumerClient[k.index][name] = make(map[string]*kafka.Consumer)
		}
		k.ConsumerClient[k.index][name][group] = consumer

		return consumer, nil
	}
	return nil, fmt.Errorf("kafka client  GetConsumerClient is error")
}

func (k *KafkaClientManager) UpdateLoadK8sConfigMap(configMapName, env string, hookName string) error {
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
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()

	oldIndex := k.index
	newIndex := (k.index + 1) % 2

	// 初始化新配置
	if k.Config[newIndex] == nil {
		k.Config[newIndex] = make(map[string]*config_manage.KafkaConfig)
	} else {
		// 清空旧配置
		for key := range k.Config[newIndex] {
			delete(k.Config[newIndex], key)
		}
	}

	// 加载新配置
	for _, v := range *e.Config {
		k.Config[newIndex][v.Name] = &v
	}

	// 切换到新配置
	k.index = newIndex

	// 异步关闭旧连接，避免阻塞和内存泄漏
	if oldIndex >= 0 {
		go k.cleanupOldConnections(oldIndex)
	}
	return nil
}
func (k *KafkaClientManager) buildProducer(conf *config_manage.KafkaConfig) (*kafka.Producer, error) {
	// 创建生产者配置
	config := &kafka.ConfigMap{
		"bootstrap.servers": conf.Producer.Broker,
		"client.id":         conf.Producer.Producer,
		"acks":              "all",
	}
	if conf.Producer.Username != EmptyString && conf.Producer.Password != EmptyString {
		config.SetKey("sasl.username", conf.Producer.Username)
		config.SetKey("sasl.password", conf.Producer.Password)
		config.SetKey("security.protocol", conf.Producer.Protocol)
		config.SetKey("sasl.mechanism", conf.Producer.Mechanism)
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %s", err)
	}

	return p, nil

}

func (k *KafkaClientManager) buildConsumer(conf *config_manage.KafkaConfig, group string) (*kafka.Consumer, error) {
	// 创建消费者配置
	config := &kafka.ConfigMap{
		"bootstrap.servers":  conf.Customer.Broker,
		"group.id":           group,
		"auto.offset.reset":  conf.Customer.Reset,
		"enable.auto.commit": conf.Customer.AutoCommit,
	}
	if conf.Customer.Username != EmptyString && conf.Customer.Password != EmptyString {
		config.SetKey("sasl.username", conf.Customer.Username)
		config.SetKey("sasl.password", conf.Customer.Password)
		config.SetKey("security.protocol", conf.Customer.Protocol)
		config.SetKey("sasl.mechanism", conf.Customer.Mechanism)
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}
	return c, nil
}

// 异步清理旧连接，防止内存泄漏
func (k *KafkaClientManager) cleanupOldConnections(oldIndex int) {
	// 关闭旧的 Producer
	if k.ProducerClient[oldIndex] != nil {
		for name, producer := range k.ProducerClient[oldIndex] {
			if producer != nil {
				producer.Close()
				fmt.Printf("Closed old producer: %s\n", name)
			}
		}
		// 清理 map
		k.rwMutex.Lock()
		k.ProducerClient[oldIndex] = make(map[string]*kafka.Producer)
		k.rwMutex.Unlock()
	}

	// 关闭旧的 Consumer
	if k.ConsumerClient[oldIndex] != nil {
		for name, consumerGroup := range k.ConsumerClient[oldIndex] {
			for group, consumer := range consumerGroup {
				if consumer != nil {
					consumer.Close()
					fmt.Printf("Closed old consumer: %s-%s\n", name, group)
				}
			}
		}
		// 清理 map
		k.rwMutex.Lock()
		k.ConsumerClient[oldIndex] = make(map[string]map[string]*kafka.Consumer)
		k.rwMutex.Unlock()
	}
}

// 关闭指定的 Producer
func (k *KafkaClientManager) CloseProducer(name string) error {
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()

	if k.index >= 0 && k.ProducerClient[k.index] != nil {
		if producer, exists := k.ProducerClient[k.index][name]; exists && producer != nil {
			producer.Close()
			delete(k.ProducerClient[k.index], name)
			fmt.Printf("Closed producer: %s\n", name)
			return nil
		}
	}
	return fmt.Errorf("producer %s not found", name)
}

// 关闭指定的 Consumer
func (k *KafkaClientManager) CloseConsumer(name string, group string) error {
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()

	if k.index >= 0 && k.ConsumerClient[k.index] != nil {
		if consumerGroup, exists := k.ConsumerClient[k.index][name]; exists && consumerGroup != nil {
			if consumer, ok := consumerGroup[group]; ok && consumer != nil {
				consumer.Close()
				delete(consumerGroup, group)
				fmt.Printf("Closed consumer: %s-%s\n", name, group)

				// 如果该 name 下没有其他 group 了，删除整个 name 项
				if len(consumerGroup) == 0 {
					delete(k.ConsumerClient[k.index], name)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("consumer %s-%s not found", name, group)
}

// 修改后的全量关闭方法（重命名为 CloseAll）
func (k *KafkaClientManager) CloseAll() {
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()

	for i := 0; i < 2; i++ {
		// 关闭所有 Producer
		if k.ProducerClient[i] != nil {
			for name, producer := range k.ProducerClient[i] {
				if producer != nil {
					producer.Close()
					fmt.Printf("Closed producer: %s\n", name)
				}
			}
			k.ProducerClient[i] = make(map[string]*kafka.Producer)
		}

		// 关闭所有 Consumer
		if k.ConsumerClient[i] != nil {
			for name, consumerGroup := range k.ConsumerClient[i] {
				for group, consumer := range consumerGroup {
					if consumer != nil {
						consumer.Close()
						fmt.Printf("Closed consumer: %s-%s\n", name, group)
					}
				}
			}
			k.ConsumerClient[i] = make(map[string]map[string]*kafka.Consumer)
		}
	}
}

// 保持原有的 Close 方法用于向后兼容，但建议使用 CloseAll
func (k *KafkaClientManager) Close() {
	k.CloseAll()
}

// func (k *KafkaClientManager) buildProducer(conf *config_manage.KafkaConfig) (*kafka.Writer, error) {
// 	// 创建生产者配置
// 	mechanism := plain.Mechanism{
// 		Username: conf.Producer.Username,
// 		Password: conf.Producer.Password,
// 	}
// 	dialer := &kafka.Dialer{
// 		SASLMechanism: mechanism,
// 		ClientID:      conf.Producer.Producer,
// 		TLS:           &tls.Config{},
// 	}
// 	p := kafka.NewWriter(kafka.WriterConfig{
// 		Brokers:  []string{conf.Producer.Broker},
// 		Balancer: &kafka.LeastBytes{},
// 		Dialer:   dialer,
// 	})
// 	// p, err := kafka.NewWriter(&kafka.WriterConfig{
// 	// 	Brokers: []string{conf.Producer.Broker},
// 	// 	Dialer:  dialer,
// 	// 	// "client.id":         conf.Producer.Producer,
// 	// 	// "acks":              "all",
// 	// 	// "security.protocol": conf.Producer.Protocol,
// 	// 	// "sasl.mechanism":    conf.Producer.Mechanism,
// 	// 	// "sasl.username":     conf.Producer.Username,
// 	// 	// "sasl.password":     conf.Producer.Password,
// 	// })
// 	return p, nil

// }

// func (k *KafkaClientManager) buildConsumer(conf *config_manage.KafkaConfig) (*kafka.Reader, error) {
// 	// 创建消费者配置
// 	// c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 	// 	"bootstrap.servers": conf.Customer.Broker,
// 	// 	"group.id":          conf.Customer.Group,
// 	// 	"auto.offset.reset": conf.Customer.Reset,
// 	// 	"security.protocol": conf.Customer.Protocol,
// 	// 	"sasl.mechanism":    conf.Customer.Mechanism,
// 	// 	"sasl.username":     conf.Customer.Username,
// 	// 	"sasl.password":     conf.Customer.Password,
// 	// })
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to create consumer: %s", err)
// 	// }
// 	mechanism := plain.Mechanism{
// 		Username: conf.Customer.Username,
// 		Password: conf.Customer.Password,
// 	}

// 	dialer := &kafka.Dialer{
// 		SASLMechanism: mechanism,
// 		TLS:           &tls.Config{},
// 	}

// 	c := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: []string{conf.Customer.Broker},
// 		GroupID: conf.Customer.Group,
// 		Dialer:  dialer,
// 	})
// 	return c, nil
// }
