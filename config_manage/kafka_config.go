package config_manage

import (
	"fmt"

	"github.com/blueturbo-ad/go-utils/environment"
	"gopkg.in/yaml.v3"
)

/*
curused: Dev
Dev:
  version: 1
  kafka_conf:
    - name: event_kafka
      producer:
        broker: bootstrap.kafka-cluster.asia-east1.managedkafka.fine-acronym-336109.cloud.goog:9092
        topic: event_topic
        producer: "event_producer"
        username: imagen-srvact@fine-acronym-336109.iam.gserviceaccount.com
        password: /home/imagen-srvact/imagen-srvact.json
        protocol: SASL_SSL
        mechanism: PLAIN
      customer:
        broker: bootstrap.kafka-cluster.asia-east1.managedkafka.fine-acronym-336109.cloud.goog:9092
        topic: event_topic
        group: event_group
        username: imagen-srvact@fine-acronym-336109.iam.gserviceaccount.com
        password: /home/imagen-srvact/imagen-srvact.json
        protocol: SASL_SSL
        mechanism: PLAIN
        reset: earliest
Pro:
  version: 1
  kafka_conf:
    - name: event_kafka
      producer:
        broker: bootstrap.kafka-cluster.asia-east1.managedkafka.fine-acronym-336109.cloud.goog:9092
        topic: event_topic
        producer: "event_producer"
        username: imagen-srvact@fine-acronym-336109.iam.gserviceaccount.com
        password: /home/imagen-srvact/imagen-srvact.json
        protocol: SASL_SSL
        mechanism: PLAIN
      customer:
        broker: bootstrap.kafka-cluster.asia-east1.managedkafka.fine-acronym-336109.cloud.goog:9092
        topic: event_topic
        group: event_group
        username: imagen-srvact@fine-acronym-336109.iam.gserviceaccount.com
        password: /home/imagen-srvact/imagen-srvact.json
        protocol: SASL_SSL
        mechanism: PLAIN
        reset: earliest
Pre: {}
Test: {}
*/

type KafkaConfig struct {
	Name     string `yaml:"name"`
	Producer struct {
		Broker    string `yaml:"broker"`
		Producer  string `yaml:"producer"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Protocol  string `yaml:"protocol"`
		Mechanism string `yaml:"mechanism"`
	} `yaml:"producer"`
	Customer struct {
		Broker    string `yaml:"broker"`
		Group     string `yaml:"group"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Protocol  string `yaml:"protocol"`
		Mechanism string `yaml:"mechanism"`
		Reset     string `yaml:"reset"`
	} `yaml:"customer"`
}

type KafkaConfigManage struct {
	Config  *[]KafkaConfig `yaml:"kafka_conf"`
	Version string         `yaml:"version"`
}

func (r *KafkaConfigManage) LoadK8sConfigMap(configMapName, env string) error {
	var c = new(ManagerConfig)
	namespace := environment.GetSingleton().GetNamespace()
	info, err := c.LoadK8sConfigMap(namespace, configMapName, env)
	if err != nil {
		return err
	}
	if (*info) == nil {
		return fmt.Errorf("info is nil，")
	}
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (r *KafkaConfigManage) LoadConfig(filePath string, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadFileConfig(filePath, env)
	if err != nil {
		return err
	}
	//fmt.Println(info)
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}

func (l *KafkaConfigManage) LoadMemoryConfig(buf []byte, env string) error {
	var c = new(ManagerConfig)
	info, err := c.LoadMemoryConfig(buf, env)
	if err != nil {
		return err
	}
	// 解析 YAML 数据
	inmap := (*info).(map[string]interface{})
	data, err := yaml.Marshal(inmap)
	if err != nil {
		return fmt.Errorf("failed to marshal inmap: %v", err)
	}
	err = yaml.Unmarshal(data, &l.Config)
	if err != nil {
		return fmt.Errorf(ErroryamlNotfound, err)
	}
	return nil
}
